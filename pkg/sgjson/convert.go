package sgjson

import (
	"github.com/aws/aws-sdk-go/service/ec2"
)

// GroupFromAWS Convert an AWS returned SecurityGroup stucture to our local structure
func GroupFromAWS(awsGroup ec2.SecurityGroup) (LocalSecurityGroup, error) {
	var group LocalSecurityGroup
	group.Ingress = make(map[string][]PortRangeEntry)

	group.Description = *awsGroup.Description
	group.GroupName = *awsGroup.GroupName
	group.GroupID = *awsGroup.GroupId
	group.OwnerID = *awsGroup.OwnerId
	group.VpcID = *awsGroup.VpcId

	for _, permission := range awsGroup.IpPermissions {
		fromPort := *permission.FromPort
		toPort := *permission.ToPort
		protocol := *permission.IpProtocol

		for _, iprange := range permission.IpRanges {
			cidrIP := *iprange.CidrIp
			description := *iprange.Description

			// Do we have an entry for that description yet ?
			if group.Ingress[description] == nil {
				// Create the slice
				group.Ingress[description] = make([]PortRangeEntry, 1, 10)

				var entry PortRangeEntry
				entry.FromPort = fromPort
				entry.ToPort = toPort
				entry.IPProtocol = protocol

				entry.IPRanges = make([]string, 1, 10)
				entry.IPRanges[0] = cidrIP

				group.Ingress[description][0] = entry
			} else {
				// Existing PortRangeEntry ?
				seen := false

				for index, iprange := range group.Ingress[description] {
					if iprange.FromPort == fromPort && iprange.ToPort == toPort && iprange.IPProtocol == protocol {
						// Add to the existing range
						iprange.IPRanges = append(iprange.IPRanges, cidrIP)
						group.Ingress[description][index] = iprange
						seen = true
					}
				}

				if seen == false {
					var entry PortRangeEntry
					entry.FromPort = fromPort
					entry.ToPort = toPort
					entry.IPProtocol = protocol

					entry.IPRanges = make([]string, 1, 10)
					entry.IPRanges[0] = cidrIP

					group.Ingress[description] = append(group.Ingress[description], entry)
				}
			}
		}

	}

	return group, nil
}
