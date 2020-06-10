package sgjson

// GroupEntry specifies a user/group id pair for a referenced Security Group
type GroupEntry struct {
	GroupID string `json:"groupId"`
	UserID  string `json:"userId"`
}

// PortRangeEntry is a set of ip's and groups that are allowed to access a port range
type PortRangeEntry struct {
	FromPort   int64        `json:"fromPort"`
	ToPort     int64        `json:"toPort"`
	IPProtocol string       `json:"ipProtocol"`
	IPRanges   []string     `json:"ipRanges"`
	GroupPairs []GroupEntry `json:"groupPairs"`
}

// LocalSecurityGroup Top level definition of OUR security group structure
type LocalSecurityGroup struct {
	Description string `json:"description"`
	GroupID     string `json:"groupId"`
	GroupName   string `json:"groupName"`
	OwnerID     string `json:"ownerId"`
	VpcID       string `json:"vpcId"`

	Ingress map[string][]PortRangeEntry // Description -> []PortRangeEntry
	Egress  map[string][]PortRangeEntry
}
