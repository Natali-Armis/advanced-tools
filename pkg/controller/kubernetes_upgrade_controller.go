package controller

import (
	"advanced-tools/pkg/client"
	"advanced-tools/pkg/entity"
	"advanced-tools/pkg/vars"
)

type K8sUpgradeController struct {
	clients *client.Client
}

func (controller *K8sUpgradeController) GetASGsNodeList() ([]*entity.ASGNodeList, error) {
	asgList := []*entity.ASGNodeList{}
	clusterAsgs, err := controller.clients.AwsClient.DescribeAutoScalingGroups(vars.DEV)
	if err != nil {
		return asgList, err
	}
	for _, clusterAsg := range clusterAsgs {
		asg := &entity.ASGNodeList{
			AsgName:  *clusterAsg.AutoScalingGroupName,
			NodeList: []*entity.AsgNode{},
		}
		asgInstances, err := controller.clients.AwsClient.ListInstances(*clusterAsg.AutoScalingGroupName)
		if err != nil {
			return asgList, err
		}
		instancesDescribe, err := controller.clients.AwsClient.DescribeInstances(asgInstances)
		if err != nil {
			return asgList, err
		}
		for _, instance := range instancesDescribe {
			version, label, err := controller.clients.K8sClient.GetNodeVersionAndLabel(*instance.PrivateDnsName)
			if err != nil {
				continue
			}
			if len(asg.Label) == 0 {
				asg.Label = label
			}
			asg.NodeList = append(asg.NodeList, &entity.AsgNode{
				InstanceId:     *instance.InstanceId,
				PrivateDnsName: *instance.PrivateDnsName,
				KubeletVersion: version,
			})
		}
		asgList = append(asgList, asg)
	}
	return asgList, nil
}

func (controller *K8sUpgradeController) EditAndVerifyIngressDeployments(newVersion string) (bool, error) {
	deployments, err := controller.clients.K8sClient.GetDeploymentsByLabelSlector(map[string]string{vars.INGRESS_LABEL_SELECTOR: vars.INGRESS_LABEL_SELECTOR_VALUE})
	if err != nil {
		return false, err
	}
	for _, deployment := range deployments {
		err = controller.clients.K8sClient.EditIngressDeploymentToMatchVersionLabel(&deployment, newVersion)
		if err != nil {
			return false, err
		}
		aligned, err := controller.clients.K8sClient.VerifyIngressDeploymentToMatchVersionLabel(&deployment, newVersion)
		if err != nil || !aligned {
			return aligned, err
		}
	}
	return true, nil
}

func (controller *K8sUpgradeController) VerifyAllNodesVersionsAligned(asgList []*entity.ASGNodeList, targetVersion string) []*entity.ASGNodeList {
	nonAlignedAsgs := []*entity.ASGNodeList{}
	for _, asg := range asgList {
		nonAlignedAsg := &entity.ASGNodeList{
			AsgName:  asg.AsgName,
			Label:    asg.Label,
			NodeList: []*entity.AsgNode{},
		}
		for _, node := range asg.NodeList {
			if node.KubeletVersion != targetVersion {
				nonAlignedAsg.NodeList = append(nonAlignedAsg.NodeList, &entity.AsgNode{
					InstanceId:     node.InstanceId,
					PrivateDnsName: node.PrivateDnsName,
					KubeletVersion: node.KubeletVersion,
				})
			}
		}
		if len(nonAlignedAsg.NodeList) > 0 {
			nonAlignedAsgs = append(nonAlignedAsgs, nonAlignedAsg)
		}
	}
	return nonAlignedAsgs
}
