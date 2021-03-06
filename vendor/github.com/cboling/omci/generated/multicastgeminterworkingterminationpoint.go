/*
 * Copyright (c) 2018 - present.  Boling Consulting Solutions (bcsw.net)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
/*
 * NOTE: This file was generated, manual edits will be overwritten!
 *
 * Generated by 'goCodeGenerator.py':
 *              https://github.com/cboling/OMCI-parser/README.md
 */
package generated

import "github.com/deckarep/golang-set"

const MulticastGemInterworkingTerminationPointClassId ClassID = ClassID(281)

var multicastgeminterworkingterminationpointBME *ManagedEntityDefinition

// MulticastGemInterworkingTerminationPoint (class ID #281)
//	An instance of this ME represents a point in a G-PON ONU where a multicast service interworks
//	with the GEM layer. At this point, a multicast bit stream is reconstructed from GEM packets.
//
//	Instances of this ME are created and deleted by the OLT.
//
//	Multicast interworking GEM modes of operation
//
//	The default multicast operation of the PON is where all the multicast content streams are
//	carried in one PON layer connection (GEM port). This connection is then specified in the first
//	entry of the IPv4 or IPv6 multicast address table, as the case may be. This single entry also
//	specifies an all-inclusive IP multicast destination address (DA) range (e.g., 224.0.0.0 to
//	239.255.255.255 in the case of IPv4). The ONU then filters the traffic based on either Ethernet
//	MAC addresses or IP addresses. The associated GEM port network CTP ME specifies the GEM port-ID
//	that supports all multicast connections.
//
//	In the default multicast operation, all multicast content streams are placed in one PON layer
//	connection (GEM port). The OLT sets up a completely conventional model, a pointer from the
//	multicast GEM IW termination to a GEM port network CTP. The OLT configures the GEM port-ID of
//	the GEM port network CTP into the appropriate multicast address table attribute(s), along with
//	the other table fields that specify the range of IP multicast DAs. The ONU accepts the entire
//	multicast stream through the designated GEM port, then filters the traffic based on either the
//	Ethernet MAC address or IP DA.
//
//	An optional multicast configuration supports separate multicast streams carried over separate
//	PON layer connections, i.e., on separate GEM ports. This permits the ONU to filter multicast
//	streams at the GEM level, which is efficient in hardware, while ignoring other multicast streams
//	that may be of interest to other ONUs on the PON.
//
//	After configuring the explicit model for the first multicast GEM port, the OLT supports multiple
//	multicast GEM ports by then configuring additional entries into the multicast address table(s),
//	entries with different GEM port-IDs. The OMCI model is defined such that these ports are
//	implicitly grouped together and served by the single explicit GEM port network CTP. No
//	additional GEM network CTPs need be created or linked for the additional GEM ports.
//
//	Several multicast GEM IW TPs can exist, each linked to separate bridge ports or mappers to serve
//	different communities of interest in a complex ONU.
//
//	Discovery of multicast support
//
//	The OLT uses the multicast GEM IW TP entity as the means to discover the ONU's multicast
//	capability. This entity is mandatory if multicast is supported by the ONU. If the OLT attempts
//	to create this entity on an ONU that does not support multicast, the create command fails. The
//	create or set command also fails if the OLT attempts to exploit optional features that the ONU
//	does not support, e.g., in attempting to write a multicast address table with more than a single
//	entry or to create multiple multicast GEM IW TPs.
//
//	This ME is defined by a similarity to the unicast GEM IW TP, and a number of its attributes are
//	not meaningful in a multicast context. These attributes are set to 0 and not used, as indicated
//	in the following.
//
//	Relationships
//		An instance of this ME exists for each occurrence of transformation of GEM packets into a
//		multicast data stream.
//
//	Attributes
//		Managed Entity Id
//			Managed entity ID: This attribute uniquely identifies each instance of this ME. The value 0xFFFF
//			is reserved. (R, setbycreate) (mandatory) (2 bytes)
//
//		Gem Port Network Ctp Connectivity Pointer
//			GEM port network CTP connectivity pointer: This attribute points to an instance of the GEM port
//			network CTP that is associated with this multicast GEM IW TP. (R, W, setbycreate) (mandatory)
//			(2 bytes)
//
//		Interworking Option
//			(R, W, setbycreate) (mandatory) (1 byte)
//
//		Service Profile Pointer
//			Service profile pointer: This attribute is set to 0 and not used. For backward compatibility, it
//			may also be set to point to a MAC bridge service profile or IEEE 802.1p mapper service profile.
//			(R, W, setbycreate) (mandatory) (2 bytes)
//
//		Pptp Counter
//			PPTP counter: This attribute represents the number of instances of PPTP MEs associated with this
//			instance of the multicast GEM IW TP. This attribute conveys no information that is not available
//			elsewhere; it may be set to 0xFF and not used. (R) (optional) (1 byte)
//
//		Operational State
//			Operational state: This attribute indicates whether the ME is capable of performing its
//			function. Valid values are enabled (0) and disabled (1). (R) (optional) (1 byte)
//
//		Gal Profile Pointer
//			GAL profile pointer: This attribute is set to 0 and not used. For backward compatibility, it may
//			also be set to point to a GAL Ethernet profile. (R, W, setbycreate) (mandatory) (2 bytes)
//
//		Ipv6 Multicast Address Table
//			(R, W) (optional) (24N bytes, where N is the number of entries in the list.)
//
type MulticastGemInterworkingTerminationPoint struct {
	ManagedEntityDefinition
	Attributes AttributeValueMap
}

func init() {
	multicastgeminterworkingterminationpointBME = &ManagedEntityDefinition{
		Name:    "MulticastGemInterworkingTerminationPoint",
		ClassID: 281,
		MessageTypes: mapset.NewSetWith(
			Create,
			Delete,
			Get,
			GetNext,
			Set,
		),
		AllowedAttributeMask: 0XFE00,
		AttributeDefinitions: AttributeDefinitionMap{
			0: Uint16Field("ManagedEntityId", 0, mapset.NewSetWith(Read, SetByCreate), false, false, false, false, 0),
			1: Uint16Field("GemPortNetworkCtpConnectivityPointer", 0, mapset.NewSetWith(Read, SetByCreate, Write), false, false, false, false, 1),
			2: ByteField("InterworkingOption", 0, mapset.NewSetWith(Read, SetByCreate, Write), false, false, false, false, 2),
			3: Uint16Field("ServiceProfilePointer", 0, mapset.NewSetWith(Read, SetByCreate, Write), false, false, false, false, 3),
			4: ByteField("PptpCounter", 0, mapset.NewSetWith(Read), false, false, true, false, 4),
			5: ByteField("OperationalState", 0, mapset.NewSetWith(Read), true, false, true, false, 5),
			6: Uint16Field("GalProfilePointer", 0, mapset.NewSetWith(Read, SetByCreate, Write), false, false, false, false, 6),
			7: TableField("Ipv6MulticastAddressTable", TableInfo{nil, 24}, mapset.NewSetWith(Read, Write), false, true, false, 7),
		},
	}
}

// NewMulticastGemInterworkingTerminationPoint (class ID 281 creates the basic
// Managed Entity definition that is used to validate an ME of this type that
// is received from the wire, about to be sent on the wire.
func NewMulticastGemInterworkingTerminationPoint(params ...ParamData) (*ManagedEntity, OmciErrors) {
	return NewManagedEntity(multicastgeminterworkingterminationpointBME, params...)
}
