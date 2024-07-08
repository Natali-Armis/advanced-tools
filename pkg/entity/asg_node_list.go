package entity

type ASGNodeList struct {
	AsgName  string
	Label    string
	NodeList []*AsgNode
}

type AsgNode struct {
	InstanceId     string
	PrivateDnsName string
	KubeletVersion string
}
