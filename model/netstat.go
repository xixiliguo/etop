package model

import (
	"github.com/xixiliguo/etop/store"
)

var (
	DefaultNetStatFields = map[string][]string{
		"tcp": {"TcpActiveOpens", "TcpPassiveOpens", "TcpAttemptFails",
			"TcpEstabResets", "TcpCurrEstab", "TcpRetransSegs", "TcpInErrs", "TcpInCsumErrors"},
		"ip": {"IpInReceives", "IpInHdrErrors", "IpInAddrErrors",
			"IpForwDatagrams", "IpInUnknownProtos", "IpInDiscards",
			"IpInDelivers", "IpOutRequests", "IpOutDiscards",
			"IpOutNoRoutes", "IpReasmTimeout", "IpReasmReqds",
			"IpReasmOKs", "IpReasmFails"},
	}
)

type NetStat struct {
	IpInReceives      float64
	IpInHdrErrors     float64
	IpInAddrErrors    float64
	IpForwDatagrams   float64
	IpInUnknownProtos float64
	IpInDiscards      float64
	IpInDelivers      float64
	IpOutRequests     float64
	IpOutDiscards     float64
	IpOutNoRoutes     float64
	IpReasmTimeout    float64
	IpReasmReqds      float64
	IpReasmOKs        float64
	IpReasmFails      float64
	IpFragOKs         float64
	IpFragFails       float64
	IpFragCreates     float64

	IcmpInMsgs           float64
	IcmpInErrors         float64
	IcmpInCsumErrors     float64
	IcmpInDestUnreachs   float64
	IcmpInTimeExcds      float64
	IcmpInParmProbs      float64
	IcmpInSrcQuenchs     float64
	IcmpInRedirects      float64
	IcmpInEchos          float64
	IcmpInEchoReps       float64
	IcmpInTimestamps     float64
	IcmpInTimestampReps  float64
	IcmpInAddrMasks      float64
	IcmpInAddrMaskReps   float64
	IcmpOutMsgs          float64
	IcmpOutErrors        float64
	IcmpOutDestUnreachs  float64
	IcmpOutTimeExcds     float64
	IcmpOutParmProbs     float64
	IcmpOutSrcQuenchs    float64
	IcmpOutRedirects     float64
	IcmpOutEchos         float64
	IcmpOutEchoReps      float64
	IcmpOutTimestamps    float64
	IcmpOutTimestampReps float64
	IcmpOutAddrMasks     float64
	IcmpOutAddrMaskReps  float64
	IcmpInType3          float64
	IcmpOutType3         float64

	TcpActiveOpens  float64
	TcpPassiveOpens float64
	TcpAttemptFails float64
	TcpEstabResets  float64
	TcpCurrEstab    float64
	TcpInSegs       float64
	TcpOutSegs      float64
	TcpRetransSegs  float64
	TcpInErrs       float64
	TcpOutRsts      float64
	TcpInCsumErrors float64

	UdpInDatagrams  float64
	UdpNoPorts      float64
	UdpInErrors     float64
	UdpOutDatagrams float64
	UdpRcvbufErrors float64
	UdpSndbufErrors float64
	UdpInCsumErrors float64
	UdpIgnoredMulti float64

	UdpLiteInDatagrams  float64
	UdpLiteNoPorts      float64
	UdpLiteInErrors     float64
	UdpLiteOutDatagrams float64
	UdpLiteRcvbufErrors float64
	UdpLiteSndbufErrors float64
	UdpLiteInCsumErrors float64
	UdpLiteIgnoredMulti float64

	Ip6InReceives       float64
	Ip6InHdrErrors      float64
	Ip6InTooBigErrors   float64
	Ip6InNoRoutes       float64
	Ip6InAddrErrors     float64
	Ip6InUnknownProtos  float64
	Ip6InTruncatedPkts  float64
	Ip6InDiscards       float64
	Ip6InDelivers       float64
	Ip6OutForwDatagrams float64
	Ip6OutRequests      float64
	Ip6OutDiscards      float64
	Ip6OutNoRoutes      float64
	Ip6ReasmTimeout     float64
	Ip6ReasmReqds       float64
	Ip6ReasmOKs         float64
	Ip6ReasmFails       float64
	Ip6FragOKs          float64
	Ip6FragFails        float64
	Ip6FragCreates      float64
	Ip6InMcastPkts      float64
	Ip6OutMcastPkts     float64
	Ip6InOctets         float64
	Ip6OutOctets        float64
	Ip6InMcastOctets    float64
	Ip6OutMcastOctets   float64
	Ip6InBcastOctets    float64
	Ip6OutBcastOctets   float64
	Ip6InNoECTPkts      float64
	Ip6InECT1Pkts       float64
	Ip6InECT0Pkts       float64
	Ip6InCEPkts         float64

	Icmp6InMsgs                    float64
	Icmp6InErrors                  float64
	Icmp6OutMsgs                   float64
	Icmp6OutErrors                 float64
	Icmp6InCsumErrors              float64
	Icmp6InDestUnreachs            float64
	Icmp6InPktTooBigs              float64
	Icmp6InTimeExcds               float64
	Icmp6InParmProblems            float64
	Icmp6InEchos                   float64
	Icmp6InEchoReplies             float64
	Icmp6InGroupMembQueries        float64
	Icmp6InGroupMembResponses      float64
	Icmp6InGroupMembReductions     float64
	Icmp6InRouterSolicits          float64
	Icmp6InRouterAdvertisements    float64
	Icmp6InNeighborSolicits        float64
	Icmp6InNeighborAdvertisements  float64
	Icmp6InRedirects               float64
	Icmp6InMLDv2Reports            float64
	Icmp6OutDestUnreachs           float64
	Icmp6OutPktTooBigs             float64
	Icmp6OutTimeExcds              float64
	Icmp6OutParmProblems           float64
	Icmp6OutEchos                  float64
	Icmp6OutEchoReplies            float64
	Icmp6OutGroupMembQueries       float64
	Icmp6OutGroupMembResponses     float64
	Icmp6OutGroupMembReductions    float64
	Icmp6OutRouterSolicits         float64
	Icmp6OutRouterAdvertisements   float64
	Icmp6OutNeighborSolicits       float64
	Icmp6OutNeighborAdvertisements float64
	Icmp6OutRedirects              float64
	Icmp6OutMLDv2Reports           float64
	Icmp6InType1                   float64
	Icmp6InType134                 float64
	Icmp6InType135                 float64
	Icmp6InType136                 float64
	Icmp6InType143                 float64
	Icmp6OutType133                float64
	Icmp6OutType135                float64
	Icmp6OutType136                float64
	Icmp6OutType143                float64

	Udp6InDatagrams  float64
	Udp6NoPorts      float64
	Udp6InErrors     float64
	Udp6OutDatagrams float64
	Udp6RcvbufErrors float64
	Udp6SndbufErrors float64
	Udp6InCsumErrors float64
	Udp6IgnoredMulti float64

	UdpLite6InDatagrams  float64
	UdpLite6NoPorts      float64
	UdpLite6InErrors     float64
	UdpLite6OutDatagrams float64
	UdpLite6RcvbufErrors float64
	UdpLite6SndbufErrors float64
	UdpLite6InCsumErrors float64

	TcpExtSyncookiesSent            float64
	TcpExtSyncookiesRecv            float64
	TcpExtSyncookiesFailed          float64
	TcpExtEmbryonicRsts             float64
	TcpExtPruneCalled               float64
	TcpExtRcvPruned                 float64
	TcpExtOfoPruned                 float64
	TcpExtOutOfWindowIcmps          float64
	TcpExtLockDroppedIcmps          float64
	TcpExtArpFilter                 float64
	TcpExtTW                        float64
	TcpExtTWRecycled                float64
	TcpExtTWKilled                  float64
	TcpExtPAWSActive                float64
	TcpExtPAWSEstab                 float64
	TcpExtDelayedACKs               float64
	TcpExtDelayedACKLocked          float64
	TcpExtDelayedACKLost            float64
	TcpExtListenOverflows           float64
	TcpExtListenDrops               float64
	TcpExtTCPHPHits                 float64
	TcpExtTCPPureAcks               float64
	TcpExtTCPHPAcks                 float64
	TcpExtTCPRenoRecovery           float64
	TcpExtTCPSackRecovery           float64
	TcpExtTCPSACKReneging           float64
	TcpExtTCPSACKReorder            float64
	TcpExtTCPRenoReorder            float64
	TcpExtTCPTSReorder              float64
	TcpExtTCPFullUndo               float64
	TcpExtTCPPartialUndo            float64
	TcpExtTCPDSACKUndo              float64
	TcpExtTCPLossUndo               float64
	TcpExtTCPLostRetransmit         float64
	TcpExtTCPRenoFailures           float64
	TcpExtTCPSackFailures           float64
	TcpExtTCPLossFailures           float64
	TcpExtTCPFastRetrans            float64
	TcpExtTCPSlowStartRetrans       float64
	TcpExtTCPTimeouts               float64
	TcpExtTCPLossProbes             float64
	TcpExtTCPLossProbeRecovery      float64
	TcpExtTCPRenoRecoveryFail       float64
	TcpExtTCPSackRecoveryFail       float64
	TcpExtTCPRcvCollapsed           float64
	TcpExtTCPDSACKOldSent           float64
	TcpExtTCPDSACKOfoSent           float64
	TcpExtTCPDSACKRecv              float64
	TcpExtTCPDSACKOfoRecv           float64
	TcpExtTCPAbortOnData            float64
	TcpExtTCPAbortOnClose           float64
	TcpExtTCPAbortOnMemory          float64
	TcpExtTCPAbortOnTimeout         float64
	TcpExtTCPAbortOnLinger          float64
	TcpExtTCPAbortFailed            float64
	TcpExtTCPMemoryPressures        float64
	TcpExtTCPMemoryPressuresChrono  float64
	TcpExtTCPSACKDiscard            float64
	TcpExtTCPDSACKIgnoredOld        float64
	TcpExtTCPDSACKIgnoredNoUndo     float64
	TcpExtTCPSpuriousRTOs           float64
	TcpExtTCPMD5NotFound            float64
	TcpExtTCPMD5Unexpected          float64
	TcpExtTCPMD5Failure             float64
	TcpExtTCPSackShifted            float64
	TcpExtTCPSackMerged             float64
	TcpExtTCPSackShiftFallback      float64
	TcpExtTCPBacklogDrop            float64
	TcpExtPFMemallocDrop            float64
	TcpExtTCPMinTTLDrop             float64
	TcpExtTCPDeferAcceptDrop        float64
	TcpExtIPReversePathFilter       float64
	TcpExtTCPTimeWaitOverflow       float64
	TcpExtTCPReqQFullDoCookies      float64
	TcpExtTCPReqQFullDrop           float64
	TcpExtTCPRetransFail            float64
	TcpExtTCPRcvCoalesce            float64
	TcpExtTCPRcvQDrop               float64
	TcpExtTCPOFOQueue               float64
	TcpExtTCPOFODrop                float64
	TcpExtTCPOFOMerge               float64
	TcpExtTCPChallengeACK           float64
	TcpExtTCPSYNChallenge           float64
	TcpExtTCPFastOpenActive         float64
	TcpExtTCPFastOpenActiveFail     float64
	TcpExtTCPFastOpenPassive        float64
	TcpExtTCPFastOpenPassiveFail    float64
	TcpExtTCPFastOpenListenOverflow float64
	TcpExtTCPFastOpenCookieReqd     float64
	TcpExtTCPFastOpenBlackhole      float64
	TcpExtTCPSpuriousRtxHostQueues  float64
	TcpExtBusyPollRxPackets         float64
	TcpExtTCPAutoCorking            float64
	TcpExtTCPFromZeroWindowAdv      float64
	TcpExtTCPToZeroWindowAdv        float64
	TcpExtTCPWantZeroWindowAdv      float64
	TcpExtTCPSynRetrans             float64
	TcpExtTCPOrigDataSent           float64
	TcpExtTCPHystartTrainDetect     float64
	TcpExtTCPHystartTrainCwnd       float64
	TcpExtTCPHystartDelayDetect     float64
	TcpExtTCPHystartDelayCwnd       float64
	TcpExtTCPACKSkippedSynRecv      float64
	TcpExtTCPACKSkippedPAWS         float64
	TcpExtTCPACKSkippedSeq          float64
	TcpExtTCPACKSkippedFinWait2     float64
	TcpExtTCPACKSkippedTimeWait     float64
	TcpExtTCPACKSkippedChallenge    float64
	TcpExtTCPWinProbe               float64
	TcpExtTCPKeepAlive              float64
	TcpExtTCPMTUPFail               float64
	TcpExtTCPMTUPSuccess            float64
	TcpExtTCPWqueueTooBig           float64

	IpExtInNoRoutes      float64
	IpExtInTruncatedPkts float64
	IpExtInMcastPkts     float64
	IpExtOutMcastPkts    float64
	IpExtInBcastPkts     float64
	IpExtOutBcastPkts    float64
	IpExtInOctets        float64
	IpExtOutOctets       float64
	IpExtInMcastOctets   float64
	IpExtOutMcastOctets  float64
	IpExtInBcastOctets   float64
	IpExtOutBcastOctets  float64
	IpExtInCsumErrors    float64
	IpExtInNoECTPkts     float64
	IpExtInECT1Pkts      float64
	IpExtInECT0Pkts      float64
	IpExtInCEPkts        float64
	IpExtReasmOverlaps   float64
}

func (n *NetStat) DefaultConfig(field string) Field {
	cfg := Field{}
	switch field {
	case "IpInReceives":
		cfg = Field{"IpInReceives", Raw, 0, "", 10, false}
	case "IpInHdrErrors":
		cfg = Field{"IpInHdrErrors", Raw, 0, "", 10, false}
	case "IpInAddrErrors":
		cfg = Field{"IpInAddrErrors", Raw, 0, "", 10, false}
	case "IpForwDatagrams":
		cfg = Field{"IpForwDatagrams", Raw, 0, "", 10, false}
	case "IpInUnknownProtos":
		cfg = Field{"IpInUnknownProtos", Raw, 0, "", 10, false}
	case "IpInDiscards":
		cfg = Field{"IpInDiscards", Raw, 0, "", 10, false}
	case "IpInDelivers":
		cfg = Field{"IpInDelivers", Raw, 0, "", 10, false}
	case "IpOutRequests":
		cfg = Field{"IpOutRequests", Raw, 0, "", 10, false}
	case "IpOutDiscards":
		cfg = Field{"IpOutDiscards", Raw, 0, "", 10, false}
	case "IpOutNoRoutes":
		cfg = Field{"IpOutNoRoutes", Raw, 0, "", 10, false}
	case "IpReasmTimeout":
		cfg = Field{"IpReasmTimeout", Raw, 0, "", 10, false}
	case "IpReasmReqds":
		cfg = Field{"IpReasmReqds", Raw, 0, "", 10, false}
	case "IpReasmOKs":
		cfg = Field{"IpReasmOKs", Raw, 0, "", 10, false}
	case "IpReasmFails":
		cfg = Field{"IpReasmFails", Raw, 0, "", 10, false}
	case "IpFragOKs":
		cfg = Field{"IpFragOKs", Raw, 0, "", 10, false}
	case "IpFragFails":
		cfg = Field{"IpFragFails", Raw, 0, "", 10, false}
	case "IpFragCreates":
		cfg = Field{"IpFragCreates", Raw, 0, "", 10, false}
	case "IcmpInMsgs":
		cfg = Field{"IcmpInMsgs", Raw, 0, "", 10, false}
	case "IcmpInErrors":
		cfg = Field{"IcmpInErrors", Raw, 0, "", 10, false}
	case "IcmpInCsumErrors":
		cfg = Field{"IcmpInCsumErrors", Raw, 0, "", 10, false}
	case "IcmpInDestUnreachs":
		cfg = Field{"IcmpInDestUnreachs", Raw, 0, "", 10, false}
	case "IcmpInTimeExcds":
		cfg = Field{"IcmpInTimeExcds", Raw, 0, "", 10, false}
	case "IcmpInParmProbs":
		cfg = Field{"IcmpInParmProbs", Raw, 0, "", 10, false}
	case "IcmpInSrcQuenchs":
		cfg = Field{"IcmpInSrcQuenchs", Raw, 0, "", 10, false}
	case "IcmpInRedirects":
		cfg = Field{"IcmpInRedirects", Raw, 0, "", 10, false}
	case "IcmpInEchos":
		cfg = Field{"IcmpInEchos", Raw, 0, "", 10, false}
	case "IcmpInEchoReps":
		cfg = Field{"IcmpInEchoReps", Raw, 0, "", 10, false}
	case "IcmpInTimestamps":
		cfg = Field{"IcmpInTimestamps", Raw, 0, "", 10, false}
	case "IcmpInTimestampReps":
		cfg = Field{"IcmpInTimestampReps", Raw, 0, "", 10, false}
	case "IcmpInAddrMasks":
		cfg = Field{"IcmpInAddrMasks", Raw, 0, "", 10, false}
	case "IcmpInAddrMaskReps":
		cfg = Field{"IcmpInAddrMaskReps", Raw, 0, "", 10, false}
	case "IcmpOutMsgs":
		cfg = Field{"IcmpOutMsgs", Raw, 0, "", 10, false}
	case "IcmpOutErrors":
		cfg = Field{"IcmpOutErrors", Raw, 0, "", 10, false}
	case "IcmpOutDestUnreachs":
		cfg = Field{"IcmpOutDestUnreachs", Raw, 0, "", 10, false}
	case "IcmpOutTimeExcds":
		cfg = Field{"IcmpOutTimeExcds", Raw, 0, "", 10, false}
	case "IcmpOutParmProbs":
		cfg = Field{"IcmpOutParmProbs", Raw, 0, "", 10, false}
	case "IcmpOutSrcQuenchs":
		cfg = Field{"IcmpOutSrcQuenchs", Raw, 0, "", 10, false}
	case "IcmpOutRedirects":
		cfg = Field{"IcmpOutRedirects", Raw, 0, "", 10, false}
	case "IcmpOutEchos":
		cfg = Field{"IcmpOutEchos", Raw, 0, "", 10, false}
	case "IcmpOutEchoReps":
		cfg = Field{"IcmpOutEchoReps", Raw, 0, "", 10, false}
	case "IcmpOutTimestamps":
		cfg = Field{"IcmpOutTimestamps", Raw, 0, "", 10, false}
	case "IcmpOutTimestampReps":
		cfg = Field{"IcmpOutTimestampReps", Raw, 0, "", 10, false}
	case "IcmpOutAddrMasks":
		cfg = Field{"IcmpOutAddrMasks", Raw, 0, "", 10, false}
	case "IcmpOutAddrMaskReps":
		cfg = Field{"IcmpOutAddrMaskReps", Raw, 0, "", 10, false}
	case "IcmpInType3":
		cfg = Field{"IcmpInType3", Raw, 0, "", 10, false}
	case "IcmpOutType3":
		cfg = Field{"IcmpOutType3", Raw, 0, "", 10, false}
	case "TcpActiveOpens":
		cfg = Field{"TcpActiveOpens", Raw, 0, "", 10, false}
	case "TcpPassiveOpens":
		cfg = Field{"TcpPassiveOpens", Raw, 0, "", 10, false}
	case "TcpAttemptFails":
		cfg = Field{"TcpAttemptFails", Raw, 0, "", 10, false}
	case "TcpEstabResets":
		cfg = Field{"TcpEstabResets", Raw, 0, "", 10, false}
	case "TcpCurrEstab":
		cfg = Field{"TcpCurrEstab", Raw, 0, "", 10, false}
	case "TcpInSegs":
		cfg = Field{"TcpInSegs", Raw, 0, "", 10, false}
	case "TcpOutSegs":
		cfg = Field{"TcpOutSegs", Raw, 0, "", 10, false}
	case "TcpRetransSegs":
		cfg = Field{"TcpRetransSegs", Raw, 0, "", 10, false}
	case "TcpInErrs":
		cfg = Field{"TcpInErrs", Raw, 0, "", 10, false}
	case "TcpOutRsts":
		cfg = Field{"TcpOutRsts", Raw, 0, "", 10, false}
	case "TcpInCsumErrors":
		cfg = Field{"TcpInCsumErrors", Raw, 0, "", 10, false}
	case "UdpInDatagrams":
		cfg = Field{"UdpInDatagrams", Raw, 0, "", 10, false}
	case "UdpNoPorts":
		cfg = Field{"UdpNoPorts", Raw, 0, "", 10, false}
	case "UdpInErrors":
		cfg = Field{"UdpInErrors", Raw, 0, "", 10, false}
	case "UdpOutDatagrams":
		cfg = Field{"UdpOutDatagrams", Raw, 0, "", 10, false}
	case "UdpRcvbufErrors":
		cfg = Field{"UdpRcvbufErrors", Raw, 0, "", 10, false}
	case "UdpSndbufErrors":
		cfg = Field{"UdpSndbufErrors", Raw, 0, "", 10, false}
	case "UdpInCsumErrors":
		cfg = Field{"UdpInCsumErrors", Raw, 0, "", 10, false}
	case "UdpIgnoredMulti":
		cfg = Field{"UdpIgnoredMulti", Raw, 0, "", 10, false}
	case "UdpLiteInDatagrams":
		cfg = Field{"UdpLiteInDatagrams", Raw, 0, "", 10, false}
	case "UdpLiteNoPorts":
		cfg = Field{"UdpLiteNoPorts", Raw, 0, "", 10, false}
	case "UdpLiteInErrors":
		cfg = Field{"UdpLiteInErrors", Raw, 0, "", 10, false}
	case "UdpLiteOutDatagrams":
		cfg = Field{"UdpLiteOutDatagrams", Raw, 0, "", 10, false}
	case "UdpLiteRcvbufErrors":
		cfg = Field{"UdpLiteRcvbufErrors", Raw, 0, "", 10, false}
	case "UdpLiteSndbufErrors":
		cfg = Field{"UdpLiteSndbufErrors", Raw, 0, "", 10, false}
	case "UdpLiteInCsumErrors":
		cfg = Field{"UdpLiteInCsumErrors", Raw, 0, "", 10, false}
	case "UdpLiteIgnoredMulti":
		cfg = Field{"UdpLiteIgnoredMulti", Raw, 0, "", 10, false}
	case "Ip6InReceives":
		cfg = Field{"Ip6InReceives", Raw, 0, "", 10, false}
	case "Ip6InHdrErrors":
		cfg = Field{"Ip6InHdrErrors", Raw, 0, "", 10, false}
	case "Ip6InTooBigErrors":
		cfg = Field{"Ip6InTooBigErrors", Raw, 0, "", 10, false}
	case "Ip6InNoRoutes":
		cfg = Field{"Ip6InNoRoutes", Raw, 0, "", 10, false}
	case "Ip6InAddrErrors":
		cfg = Field{"Ip6InAddrErrors", Raw, 0, "", 10, false}
	case "Ip6InUnknownProtos":
		cfg = Field{"Ip6InUnknownProtos", Raw, 0, "", 10, false}
	case "Ip6InTruncatedPkts":
		cfg = Field{"Ip6InTruncatedPkts", Raw, 0, "", 10, false}
	case "Ip6InDiscards":
		cfg = Field{"Ip6InDiscards", Raw, 0, "", 10, false}
	case "Ip6InDelivers":
		cfg = Field{"Ip6InDelivers", Raw, 0, "", 10, false}
	case "Ip6OutForwDatagrams":
		cfg = Field{"Ip6OutForwDatagrams", Raw, 0, "", 10, false}
	case "Ip6OutRequests":
		cfg = Field{"Ip6OutRequests", Raw, 0, "", 10, false}
	case "Ip6OutDiscards":
		cfg = Field{"Ip6OutDiscards", Raw, 0, "", 10, false}
	case "Ip6OutNoRoutes":
		cfg = Field{"Ip6OutNoRoutes", Raw, 0, "", 10, false}
	case "Ip6ReasmTimeout":
		cfg = Field{"Ip6ReasmTimeout", Raw, 0, "", 10, false}
	case "Ip6ReasmReqds":
		cfg = Field{"Ip6ReasmReqds", Raw, 0, "", 10, false}
	case "Ip6ReasmOKs":
		cfg = Field{"Ip6ReasmOKs", Raw, 0, "", 10, false}
	case "Ip6ReasmFails":
		cfg = Field{"Ip6ReasmFails", Raw, 0, "", 10, false}
	case "Ip6FragOKs":
		cfg = Field{"Ip6FragOKs", Raw, 0, "", 10, false}
	case "Ip6FragFails":
		cfg = Field{"Ip6FragFails", Raw, 0, "", 10, false}
	case "Ip6FragCreates":
		cfg = Field{"Ip6FragCreates", Raw, 0, "", 10, false}
	case "Ip6InMcastPkts":
		cfg = Field{"Ip6InMcastPkts", Raw, 0, "", 10, false}
	case "Ip6OutMcastPkts":
		cfg = Field{"Ip6OutMcastPkts", Raw, 0, "", 10, false}
	case "Ip6InOctets":
		cfg = Field{"Ip6InOctets", Raw, 0, "", 10, false}
	case "Ip6OutOctets":
		cfg = Field{"Ip6OutOctets", Raw, 0, "", 10, false}
	case "Ip6InMcastOctets":
		cfg = Field{"Ip6InMcastOctets", Raw, 0, "", 10, false}
	case "Ip6OutMcastOctets":
		cfg = Field{"Ip6OutMcastOctets", Raw, 0, "", 10, false}
	case "Ip6InBcastOctets":
		cfg = Field{"Ip6InBcastOctets", Raw, 0, "", 10, false}
	case "Ip6OutBcastOctets":
		cfg = Field{"Ip6OutBcastOctets", Raw, 0, "", 10, false}
	case "Ip6InNoECTPkts":
		cfg = Field{"Ip6InNoECTPkts", Raw, 0, "", 10, false}
	case "Ip6InECT1Pkts":
		cfg = Field{"Ip6InECT1Pkts", Raw, 0, "", 10, false}
	case "Ip6InECT0Pkts":
		cfg = Field{"Ip6InECT0Pkts", Raw, 0, "", 10, false}
	case "Ip6InCEPkts":
		cfg = Field{"Ip6InCEPkts", Raw, 0, "", 10, false}
	case "Icmp6InMsgs":
		cfg = Field{"Icmp6InMsgs", Raw, 0, "", 10, false}
	case "Icmp6InErrors":
		cfg = Field{"Icmp6InErrors", Raw, 0, "", 10, false}
	case "Icmp6OutMsgs":
		cfg = Field{"Icmp6OutMsgs", Raw, 0, "", 10, false}
	case "Icmp6OutErrors":
		cfg = Field{"Icmp6OutErrors", Raw, 0, "", 10, false}
	case "Icmp6InCsumErrors":
		cfg = Field{"Icmp6InCsumErrors", Raw, 0, "", 10, false}
	case "Icmp6InDestUnreachs":
		cfg = Field{"Icmp6InDestUnreachs", Raw, 0, "", 10, false}
	case "Icmp6InPktTooBigs":
		cfg = Field{"Icmp6InPktTooBigs", Raw, 0, "", 10, false}
	case "Icmp6InTimeExcds":
		cfg = Field{"Icmp6InTimeExcds", Raw, 0, "", 10, false}
	case "Icmp6InParmProblems":
		cfg = Field{"Icmp6InParmProblems", Raw, 0, "", 10, false}
	case "Icmp6InEchos":
		cfg = Field{"Icmp6InEchos", Raw, 0, "", 10, false}
	case "Icmp6InEchoReplies":
		cfg = Field{"Icmp6InEchoReplies", Raw, 0, "", 10, false}
	case "Icmp6InGroupMembQueries":
		cfg = Field{"Icmp6InGroupMembQueries", Raw, 0, "", 10, false}
	case "Icmp6InGroupMembResponses":
		cfg = Field{"Icmp6InGroupMembResponses", Raw, 0, "", 10, false}
	case "Icmp6InGroupMembReductions":
		cfg = Field{"Icmp6InGroupMembReductions", Raw, 0, "", 10, false}
	case "Icmp6InRouterSolicits":
		cfg = Field{"Icmp6InRouterSolicits", Raw, 0, "", 10, false}
	case "Icmp6InRouterAdvertisements":
		cfg = Field{"Icmp6InRouterAdvertisements", Raw, 0, "", 10, false}
	case "Icmp6InNeighborSolicits":
		cfg = Field{"Icmp6InNeighborSolicits", Raw, 0, "", 10, false}
	case "Icmp6InNeighborAdvertisements":
		cfg = Field{"Icmp6InNeighborAdvertisements", Raw, 0, "", 10, false}
	case "Icmp6InRedirects":
		cfg = Field{"Icmp6InRedirects", Raw, 0, "", 10, false}
	case "Icmp6InMLDv2Reports":
		cfg = Field{"Icmp6InMLDv2Reports", Raw, 0, "", 10, false}
	case "Icmp6OutDestUnreachs":
		cfg = Field{"Icmp6OutDestUnreachs", Raw, 0, "", 10, false}
	case "Icmp6OutPktTooBigs":
		cfg = Field{"Icmp6OutPktTooBigs", Raw, 0, "", 10, false}
	case "Icmp6OutTimeExcds":
		cfg = Field{"Icmp6OutTimeExcds", Raw, 0, "", 10, false}
	case "Icmp6OutParmProblems":
		cfg = Field{"Icmp6OutParmProblems", Raw, 0, "", 10, false}
	case "Icmp6OutEchos":
		cfg = Field{"Icmp6OutEchos", Raw, 0, "", 10, false}
	case "Icmp6OutEchoReplies":
		cfg = Field{"Icmp6OutEchoReplies", Raw, 0, "", 10, false}
	case "Icmp6OutGroupMembQueries":
		cfg = Field{"Icmp6OutGroupMembQueries", Raw, 0, "", 10, false}
	case "Icmp6OutGroupMembResponses":
		cfg = Field{"Icmp6OutGroupMembResponses", Raw, 0, "", 10, false}
	case "Icmp6OutGroupMembReductions":
		cfg = Field{"Icmp6OutGroupMembReductions", Raw, 0, "", 10, false}
	case "Icmp6OutRouterSolicits":
		cfg = Field{"Icmp6OutRouterSolicits", Raw, 0, "", 10, false}
	case "Icmp6OutRouterAdvertisements":
		cfg = Field{"Icmp6OutRouterAdvertisements", Raw, 0, "", 10, false}
	case "Icmp6OutNeighborSolicits":
		cfg = Field{"Icmp6OutNeighborSolicits", Raw, 0, "", 10, false}
	case "Icmp6OutNeighborAdvertisements":
		cfg = Field{"Icmp6OutNeighborAdvertisements", Raw, 0, "", 10, false}
	case "Icmp6OutRedirects":
		cfg = Field{"Icmp6OutRedirects", Raw, 0, "", 10, false}
	case "Icmp6OutMLDv2Reports":
		cfg = Field{"Icmp6OutMLDv2Reports", Raw, 0, "", 10, false}
	case "Icmp6InType1":
		cfg = Field{"Icmp6InType1", Raw, 0, "", 10, false}
	case "Icmp6InType134":
		cfg = Field{"Icmp6InType134", Raw, 0, "", 10, false}
	case "Icmp6InType135":
		cfg = Field{"Icmp6InType135", Raw, 0, "", 10, false}
	case "Icmp6InType136":
		cfg = Field{"Icmp6InType136", Raw, 0, "", 10, false}
	case "Icmp6InType143":
		cfg = Field{"Icmp6InType143", Raw, 0, "", 10, false}
	case "Icmp6OutType133":
		cfg = Field{"Icmp6OutType133", Raw, 0, "", 10, false}
	case "Icmp6OutType135":
		cfg = Field{"Icmp6OutType135", Raw, 0, "", 10, false}
	case "Icmp6OutType136":
		cfg = Field{"Icmp6OutType136", Raw, 0, "", 10, false}
	case "Icmp6OutType143":
		cfg = Field{"Icmp6OutType143", Raw, 0, "", 10, false}
	case "Udp6InDatagrams":
		cfg = Field{"Udp6InDatagrams", Raw, 0, "", 10, false}
	case "Udp6NoPorts":
		cfg = Field{"Udp6NoPorts", Raw, 0, "", 10, false}
	case "Udp6InErrors":
		cfg = Field{"Udp6InErrors", Raw, 0, "", 10, false}
	case "Udp6OutDatagrams":
		cfg = Field{"Udp6OutDatagrams", Raw, 0, "", 10, false}
	case "Udp6RcvbufErrors":
		cfg = Field{"Udp6RcvbufErrors", Raw, 0, "", 10, false}
	case "Udp6SndbufErrors":
		cfg = Field{"Udp6SndbufErrors", Raw, 0, "", 10, false}
	case "Udp6InCsumErrors":
		cfg = Field{"Udp6InCsumErrors", Raw, 0, "", 10, false}
	case "Udp6IgnoredMulti":
		cfg = Field{"Udp6IgnoredMulti", Raw, 0, "", 10, false}
	case "UdpLite6InDatagrams":
		cfg = Field{"UdpLite6InDatagrams", Raw, 0, "", 10, false}
	case "UdpLite6NoPorts":
		cfg = Field{"UdpLite6NoPorts", Raw, 0, "", 10, false}
	case "UdpLite6InErrors":
		cfg = Field{"UdpLite6InErrors", Raw, 0, "", 10, false}
	case "UdpLite6OutDatagrams":
		cfg = Field{"UdpLite6OutDatagrams", Raw, 0, "", 10, false}
	case "UdpLite6RcvbufErrors":
		cfg = Field{"UdpLite6RcvbufErrors", Raw, 0, "", 10, false}
	case "UdpLite6SndbufErrors":
		cfg = Field{"UdpLite6SndbufErrors", Raw, 0, "", 10, false}
	case "UdpLite6InCsumErrors":
		cfg = Field{"UdpLite6InCsumErrors", Raw, 0, "", 10, false}
	case "TcpExtSyncookiesSent":
		cfg = Field{"TcpExtSyncookiesSent", Raw, 0, "", 10, false}
	case "TcpExtSyncookiesRecv":
		cfg = Field{"TcpExtSyncookiesRecv", Raw, 0, "", 10, false}
	case "TcpExtSyncookiesFailed":
		cfg = Field{"TcpExtSyncookiesFailed", Raw, 0, "", 10, false}
	case "TcpExtEmbryonicRsts":
		cfg = Field{"TcpExtEmbryonicRsts", Raw, 0, "", 10, false}
	case "TcpExtPruneCalled":
		cfg = Field{"TcpExtPruneCalled", Raw, 0, "", 10, false}
	case "TcpExtRcvPruned":
		cfg = Field{"TcpExtRcvPruned", Raw, 0, "", 10, false}
	case "TcpExtOfoPruned":
		cfg = Field{"TcpExtOfoPruned", Raw, 0, "", 10, false}
	case "TcpExtOutOfWindowIcmps":
		cfg = Field{"TcpExtOutOfWindowIcmps", Raw, 0, "", 10, false}
	case "TcpExtLockDroppedIcmps":
		cfg = Field{"TcpExtLockDroppedIcmps", Raw, 0, "", 10, false}
	case "TcpExtArpFilter":
		cfg = Field{"TcpExtArpFilter", Raw, 0, "", 10, false}
	case "TcpExtTW":
		cfg = Field{"TcpExtTW", Raw, 0, "", 10, false}
	case "TcpExtTWRecycled":
		cfg = Field{"TcpExtTWRecycled", Raw, 0, "", 10, false}
	case "TcpExtTWKilled":
		cfg = Field{"TcpExtTWKilled", Raw, 0, "", 10, false}
	case "TcpExtPAWSActive":
		cfg = Field{"TcpExtPAWSActive", Raw, 0, "", 10, false}
	case "TcpExtPAWSEstab":
		cfg = Field{"TcpExtPAWSEstab", Raw, 0, "", 10, false}
	case "TcpExtDelayedACKs":
		cfg = Field{"TcpExtDelayedACKs", Raw, 0, "", 10, false}
	case "TcpExtDelayedACKLocked":
		cfg = Field{"TcpExtDelayedACKLocked", Raw, 0, "", 10, false}
	case "TcpExtDelayedACKLost":
		cfg = Field{"TcpExtDelayedACKLost", Raw, 0, "", 10, false}
	case "TcpExtListenOverflows":
		cfg = Field{"TcpExtListenOverflows", Raw, 0, "", 10, false}
	case "TcpExtListenDrops":
		cfg = Field{"TcpExtListenDrops", Raw, 0, "", 10, false}
	case "TcpExtTCPHPHits":
		cfg = Field{"TcpExtTCPHPHits", Raw, 0, "", 10, false}
	case "TcpExtTCPPureAcks":
		cfg = Field{"TcpExtTCPPureAcks", Raw, 0, "", 10, false}
	case "TcpExtTCPHPAcks":
		cfg = Field{"TcpExtTCPHPAcks", Raw, 0, "", 10, false}
	case "TcpExtTCPRenoRecovery":
		cfg = Field{"TcpExtTCPRenoRecovery", Raw, 0, "", 10, false}
	case "TcpExtTCPSackRecovery":
		cfg = Field{"TcpExtTCPSackRecovery", Raw, 0, "", 10, false}
	case "TcpExtTCPSACKReneging":
		cfg = Field{"TcpExtTCPSACKReneging", Raw, 0, "", 10, false}
	case "TcpExtTCPSACKReorder":
		cfg = Field{"TcpExtTCPSACKReorder", Raw, 0, "", 10, false}
	case "TcpExtTCPRenoReorder":
		cfg = Field{"TcpExtTCPRenoReorder", Raw, 0, "", 10, false}
	case "TcpExtTCPTSReorder":
		cfg = Field{"TcpExtTCPTSReorder", Raw, 0, "", 10, false}
	case "TcpExtTCPFullUndo":
		cfg = Field{"TcpExtTCPFullUndo", Raw, 0, "", 10, false}
	case "TcpExtTCPPartialUndo":
		cfg = Field{"TcpExtTCPPartialUndo", Raw, 0, "", 10, false}
	case "TcpExtTCPDSACKUndo":
		cfg = Field{"TcpExtTCPDSACKUndo", Raw, 0, "", 10, false}
	case "TcpExtTCPLossUndo":
		cfg = Field{"TcpExtTCPLossUndo", Raw, 0, "", 10, false}
	case "TcpExtTCPLostRetransmit":
		cfg = Field{"TcpExtTCPLostRetransmit", Raw, 0, "", 10, false}
	case "TcpExtTCPRenoFailures":
		cfg = Field{"TcpExtTCPRenoFailures", Raw, 0, "", 10, false}
	case "TcpExtTCPSackFailures":
		cfg = Field{"TcpExtTCPSackFailures", Raw, 0, "", 10, false}
	case "TcpExtTCPLossFailures":
		cfg = Field{"TcpExtTCPLossFailures", Raw, 0, "", 10, false}
	case "TcpExtTCPFastRetrans":
		cfg = Field{"TcpExtTCPFastRetrans", Raw, 0, "", 10, false}
	case "TcpExtTCPSlowStartRetrans":
		cfg = Field{"TcpExtTCPSlowStartRetrans", Raw, 0, "", 10, false}
	case "TcpExtTCPTimeouts":
		cfg = Field{"TcpExtTCPTimeouts", Raw, 0, "", 10, false}
	case "TcpExtTCPLossProbes":
		cfg = Field{"TcpExtTCPLossProbes", Raw, 0, "", 10, false}
	case "TcpExtTCPLossProbeRecovery":
		cfg = Field{"TcpExtTCPLossProbeRecovery", Raw, 0, "", 10, false}
	case "TcpExtTCPRenoRecoveryFail":
		cfg = Field{"TcpExtTCPRenoRecoveryFail", Raw, 0, "", 10, false}
	case "TcpExtTCPSackRecoveryFail":
		cfg = Field{"TcpExtTCPSackRecoveryFail", Raw, 0, "", 10, false}
	case "TcpExtTCPRcvCollapsed":
		cfg = Field{"TcpExtTCPRcvCollapsed", Raw, 0, "", 10, false}
	case "TcpExtTCPDSACKOldSent":
		cfg = Field{"TcpExtTCPDSACKOldSent", Raw, 0, "", 10, false}
	case "TcpExtTCPDSACKOfoSent":
		cfg = Field{"TcpExtTCPDSACKOfoSent", Raw, 0, "", 10, false}
	case "TcpExtTCPDSACKRecv":
		cfg = Field{"TcpExtTCPDSACKRecv", Raw, 0, "", 10, false}
	case "TcpExtTCPDSACKOfoRecv":
		cfg = Field{"TcpExtTCPDSACKOfoRecv", Raw, 0, "", 10, false}
	case "TcpExtTCPAbortOnData":
		cfg = Field{"TcpExtTCPAbortOnData", Raw, 0, "", 10, false}
	case "TcpExtTCPAbortOnClose":
		cfg = Field{"TcpExtTCPAbortOnClose", Raw, 0, "", 10, false}
	case "TcpExtTCPAbortOnMemory":
		cfg = Field{"TcpExtTCPAbortOnMemory", Raw, 0, "", 10, false}
	case "TcpExtTCPAbortOnTimeout":
		cfg = Field{"TcpExtTCPAbortOnTimeout", Raw, 0, "", 10, false}
	case "TcpExtTCPAbortOnLinger":
		cfg = Field{"TcpExtTCPAbortOnLinger", Raw, 0, "", 10, false}
	case "TcpExtTCPAbortFailed":
		cfg = Field{"TcpExtTCPAbortFailed", Raw, 0, "", 10, false}
	case "TcpExtTCPMemoryPressures":
		cfg = Field{"TcpExtTCPMemoryPressures", Raw, 0, "", 10, false}
	case "TcpExtTCPMemoryPressuresChrono":
		cfg = Field{"TcpExtTCPMemoryPressuresChrono", Raw, 0, "", 10, false}
	case "TcpExtTCPSACKDiscard":
		cfg = Field{"TcpExtTCPSACKDiscard", Raw, 0, "", 10, false}
	case "TcpExtTCPDSACKIgnoredOld":
		cfg = Field{"TcpExtTCPDSACKIgnoredOld", Raw, 0, "", 10, false}
	case "TcpExtTCPDSACKIgnoredNoUndo":
		cfg = Field{"TcpExtTCPDSACKIgnoredNoUndo", Raw, 0, "", 10, false}
	case "TcpExtTCPSpuriousRTOs":
		cfg = Field{"TcpExtTCPSpuriousRTOs", Raw, 0, "", 10, false}
	case "TcpExtTCPMD5NotFound":
		cfg = Field{"TcpExtTCPMD5NotFound", Raw, 0, "", 10, false}
	case "TcpExtTCPMD5Unexpected":
		cfg = Field{"TcpExtTCPMD5Unexpected", Raw, 0, "", 10, false}
	case "TcpExtTCPMD5Failure":
		cfg = Field{"TcpExtTCPMD5Failure", Raw, 0, "", 10, false}
	case "TcpExtTCPSackShifted":
		cfg = Field{"TcpExtTCPSackShifted", Raw, 0, "", 10, false}
	case "TcpExtTCPSackMerged":
		cfg = Field{"TcpExtTCPSackMerged", Raw, 0, "", 10, false}
	case "TcpExtTCPSackShiftFallback":
		cfg = Field{"TcpExtTCPSackShiftFallback", Raw, 0, "", 10, false}
	case "TcpExtTCPBacklogDrop":
		cfg = Field{"TcpExtTCPBacklogDrop", Raw, 0, "", 10, false}
	case "TcpExtPFMemallocDrop":
		cfg = Field{"TcpExtPFMemallocDrop", Raw, 0, "", 10, false}
	case "TcpExtTCPMinTTLDrop":
		cfg = Field{"TcpExtTCPMinTTLDrop", Raw, 0, "", 10, false}
	case "TcpExtTCPDeferAcceptDrop":
		cfg = Field{"TcpExtTCPDeferAcceptDrop", Raw, 0, "", 10, false}
	case "TcpExtIPReversePathFilter":
		cfg = Field{"TcpExtIPReversePathFilter", Raw, 0, "", 10, false}
	case "TcpExtTCPTimeWaitOverflow":
		cfg = Field{"TcpExtTCPTimeWaitOverflow", Raw, 0, "", 10, false}
	case "TcpExtTCPReqQFullDoCookies":
		cfg = Field{"TcpExtTCPReqQFullDoCookies", Raw, 0, "", 10, false}
	case "TcpExtTCPReqQFullDrop":
		cfg = Field{"TcpExtTCPReqQFullDrop", Raw, 0, "", 10, false}
	case "TcpExtTCPRetransFail":
		cfg = Field{"TcpExtTCPRetransFail", Raw, 0, "", 10, false}
	case "TcpExtTCPRcvCoalesce":
		cfg = Field{"TcpExtTCPRcvCoalesce", Raw, 0, "", 10, false}
	case "TcpExtTCPRcvQDrop":
		cfg = Field{"TcpExtTCPRcvQDrop", Raw, 0, "", 10, false}
	case "TcpExtTCPOFOQueue":
		cfg = Field{"TcpExtTCPOFOQueue", Raw, 0, "", 10, false}
	case "TcpExtTCPOFODrop":
		cfg = Field{"TcpExtTCPOFODrop", Raw, 0, "", 10, false}
	case "TcpExtTCPOFOMerge":
		cfg = Field{"TcpExtTCPOFOMerge", Raw, 0, "", 10, false}
	case "TcpExtTCPChallengeACK":
		cfg = Field{"TcpExtTCPChallengeACK", Raw, 0, "", 10, false}
	case "TcpExtTCPSYNChallenge":
		cfg = Field{"TcpExtTCPSYNChallenge", Raw, 0, "", 10, false}
	case "TcpExtTCPFastOpenActive":
		cfg = Field{"TcpExtTCPFastOpenActive", Raw, 0, "", 10, false}
	case "TcpExtTCPFastOpenActiveFail":
		cfg = Field{"TcpExtTCPFastOpenActiveFail", Raw, 0, "", 10, false}
	case "TcpExtTCPFastOpenPassive":
		cfg = Field{"TcpExtTCPFastOpenPassive", Raw, 0, "", 10, false}
	case "TcpExtTCPFastOpenPassiveFail":
		cfg = Field{"TcpExtTCPFastOpenPassiveFail", Raw, 0, "", 10, false}
	case "TcpExtTCPFastOpenListenOverflow":
		cfg = Field{"TcpExtTCPFastOpenListenOverflow", Raw, 0, "", 10, false}
	case "TcpExtTCPFastOpenCookieReqd":
		cfg = Field{"TcpExtTCPFastOpenCookieReqd", Raw, 0, "", 10, false}
	case "TcpExtTCPFastOpenBlackhole":
		cfg = Field{"TcpExtTCPFastOpenBlackhole", Raw, 0, "", 10, false}
	case "TcpExtTCPSpuriousRtxHostQueues":
		cfg = Field{"TcpExtTCPSpuriousRtxHostQueues", Raw, 0, "", 10, false}
	case "TcpExtBusyPollRxPackets":
		cfg = Field{"TcpExtBusyPollRxPackets", Raw, 0, "", 10, false}
	case "TcpExtTCPAutoCorking":
		cfg = Field{"TcpExtTCPAutoCorking", Raw, 0, "", 10, false}
	case "TcpExtTCPFromZeroWindowAdv":
		cfg = Field{"TcpExtTCPFromZeroWindowAdv", Raw, 0, "", 10, false}
	case "TcpExtTCPToZeroWindowAdv":
		cfg = Field{"TcpExtTCPToZeroWindowAdv", Raw, 0, "", 10, false}
	case "TcpExtTCPWantZeroWindowAdv":
		cfg = Field{"TcpExtTCPWantZeroWindowAdv", Raw, 0, "", 10, false}
	case "TcpExtTCPSynRetrans":
		cfg = Field{"TcpExtTCPSynRetrans", Raw, 0, "", 10, false}
	case "TcpExtTCPOrigDataSent":
		cfg = Field{"TcpExtTCPOrigDataSent", Raw, 0, "", 10, false}
	case "TcpExtTCPHystartTrainDetect":
		cfg = Field{"TcpExtTCPHystartTrainDetect", Raw, 0, "", 10, false}
	case "TcpExtTCPHystartTrainCwnd":
		cfg = Field{"TcpExtTCPHystartTrainCwnd", Raw, 0, "", 10, false}
	case "TcpExtTCPHystartDelayDetect":
		cfg = Field{"TcpExtTCPHystartDelayDetect", Raw, 0, "", 10, false}
	case "TcpExtTCPHystartDelayCwnd":
		cfg = Field{"TcpExtTCPHystartDelayCwnd", Raw, 0, "", 10, false}
	case "TcpExtTCPACKSkippedSynRecv":
		cfg = Field{"TcpExtTCPACKSkippedSynRecv", Raw, 0, "", 10, false}
	case "TcpExtTCPACKSkippedPAWS":
		cfg = Field{"TcpExtTCPACKSkippedPAWS", Raw, 0, "", 10, false}
	case "TcpExtTCPACKSkippedSeq":
		cfg = Field{"TcpExtTCPACKSkippedSeq", Raw, 0, "", 10, false}
	case "TcpExtTCPACKSkippedFinWait2":
		cfg = Field{"TcpExtTCPACKSkippedFinWait2", Raw, 0, "", 10, false}
	case "TcpExtTCPACKSkippedTimeWait":
		cfg = Field{"TcpExtTCPACKSkippedTimeWait", Raw, 0, "", 10, false}
	case "TcpExtTCPACKSkippedChallenge":
		cfg = Field{"TcpExtTCPACKSkippedChallenge", Raw, 0, "", 10, false}
	case "TcpExtTCPWinProbe":
		cfg = Field{"TcpExtTCPWinProbe", Raw, 0, "", 10, false}
	case "TcpExtTCPKeepAlive":
		cfg = Field{"TcpExtTCPKeepAlive", Raw, 0, "", 10, false}
	case "TcpExtTCPMTUPFail":
		cfg = Field{"TcpExtTCPMTUPFail", Raw, 0, "", 10, false}
	case "TcpExtTCPMTUPSuccess":
		cfg = Field{"TcpExtTCPMTUPSuccess", Raw, 0, "", 10, false}
	case "TcpExtTCPWqueueTooBig":
		cfg = Field{"TcpExtTCPWqueueTooBig", Raw, 0, "", 10, false}
	case "IpExtInNoRoutes":
		cfg = Field{"IpExtInNoRoutes", Raw, 0, "", 10, false}
	case "IpExtInTruncatedPkts":
		cfg = Field{"IpExtInTruncatedPkts", Raw, 0, "", 10, false}
	case "IpExtInMcastPkts":
		cfg = Field{"IpExtInMcastPkts", Raw, 0, "", 10, false}
	case "IpExtOutMcastPkts":
		cfg = Field{"IpExtOutMcastPkts", Raw, 0, "", 10, false}
	case "IpExtInBcastPkts":
		cfg = Field{"IpExtInBcastPkts", Raw, 0, "", 10, false}
	case "IpExtOutBcastPkts":
		cfg = Field{"IpExtOutBcastPkts", Raw, 0, "", 10, false}
	case "IpExtInOctets":
		cfg = Field{"IpExtInOctets", Raw, 0, "", 10, false}
	case "IpExtOutOctets":
		cfg = Field{"IpExtOutOctets", Raw, 0, "", 10, false}
	case "IpExtInMcastOctets":
		cfg = Field{"IpExtInMcastOctets", Raw, 0, "", 10, false}
	case "IpExtOutMcastOctets":
		cfg = Field{"IpExtOutMcastOctets", Raw, 0, "", 10, false}
	case "IpExtInBcastOctets":
		cfg = Field{"IpExtInBcastOctets", Raw, 0, "", 10, false}
	case "IpExtOutBcastOctets":
		cfg = Field{"IpExtOutBcastOctets", Raw, 0, "", 10, false}
	case "IpExtInCsumErrors":
		cfg = Field{"IpExtInCsumErrors", Raw, 0, "", 10, false}
	case "IpExtInNoECTPkts":
		cfg = Field{"IpExtInNoECTPkts", Raw, 0, "", 10, false}
	case "IpExtInECT1Pkts":
		cfg = Field{"IpExtInECT1Pkts", Raw, 0, "", 10, false}
	case "IpExtInECT0Pkts":
		cfg = Field{"IpExtInECT0Pkts", Raw, 0, "", 10, false}
	case "IpExtInCEPkts":
		cfg = Field{"IpExtInCEPkts", Raw, 0, "", 10, false}
	case "IpExtReasmOverlaps":
		cfg = Field{"IpExtReasmOverlaps", Raw, 0, "", 10, false}
	}
	return cfg
}

func (n *NetStat) GetRenderValue(field string, opt FieldOpt) string {

	cfg := n.DefaultConfig(field)
	cfg.ApplyOpt(opt)
	s := ""
	switch field {
	case "IpInReceives":
		s = cfg.Render(n.IpInReceives)
	case "IpInHdrErrors":
		s = cfg.Render(n.IpInHdrErrors)
	case "IpInAddrErrors":
		s = cfg.Render(n.IpInAddrErrors)
	case "IpForwDatagrams":
		s = cfg.Render(n.IpForwDatagrams)
	case "IpInUnknownProtos":
		s = cfg.Render(n.IpInUnknownProtos)
	case "IpInDiscards":
		s = cfg.Render(n.IpInDiscards)
	case "IpInDelivers":
		s = cfg.Render(n.IpInDelivers)
	case "IpOutRequests":
		s = cfg.Render(n.IpOutRequests)
	case "IpOutDiscards":
		s = cfg.Render(n.IpOutDiscards)
	case "IpOutNoRoutes":
		s = cfg.Render(n.IpOutNoRoutes)
	case "IpReasmTimeout":
		s = cfg.Render(n.IpReasmTimeout)
	case "IpReasmReqds":
		s = cfg.Render(n.IpReasmReqds)
	case "IpReasmOKs":
		s = cfg.Render(n.IpReasmOKs)
	case "IpReasmFails":
		s = cfg.Render(n.IpReasmFails)
	case "IpFragOKs":
		s = cfg.Render(n.IpFragOKs)
	case "IpFragFails":
		s = cfg.Render(n.IpFragFails)
	case "IpFragCreates":
		s = cfg.Render(n.IpFragCreates)
	case "IcmpInMsgs":
		s = cfg.Render(n.IcmpInMsgs)
	case "IcmpInErrors":
		s = cfg.Render(n.IcmpInErrors)
	case "IcmpInCsumErrors":
		s = cfg.Render(n.IcmpInCsumErrors)
	case "IcmpInDestUnreachs":
		s = cfg.Render(n.IcmpInDestUnreachs)
	case "IcmpInTimeExcds":
		s = cfg.Render(n.IcmpInTimeExcds)
	case "IcmpInParmProbs":
		s = cfg.Render(n.IcmpInParmProbs)
	case "IcmpInSrcQuenchs":
		s = cfg.Render(n.IcmpInSrcQuenchs)
	case "IcmpInRedirects":
		s = cfg.Render(n.IcmpInRedirects)
	case "IcmpInEchos":
		s = cfg.Render(n.IcmpInEchos)
	case "IcmpInEchoReps":
		s = cfg.Render(n.IcmpInEchoReps)
	case "IcmpInTimestamps":
		s = cfg.Render(n.IcmpInTimestamps)
	case "IcmpInTimestampReps":
		s = cfg.Render(n.IcmpInTimestampReps)
	case "IcmpInAddrMasks":
		s = cfg.Render(n.IcmpInAddrMasks)
	case "IcmpInAddrMaskReps":
		s = cfg.Render(n.IcmpInAddrMaskReps)
	case "IcmpOutMsgs":
		s = cfg.Render(n.IcmpOutMsgs)
	case "IcmpOutErrors":
		s = cfg.Render(n.IcmpOutErrors)
	case "IcmpOutDestUnreachs":
		s = cfg.Render(n.IcmpOutDestUnreachs)
	case "IcmpOutTimeExcds":
		s = cfg.Render(n.IcmpOutTimeExcds)
	case "IcmpOutParmProbs":
		s = cfg.Render(n.IcmpOutParmProbs)
	case "IcmpOutSrcQuenchs":
		s = cfg.Render(n.IcmpOutSrcQuenchs)
	case "IcmpOutRedirects":
		s = cfg.Render(n.IcmpOutRedirects)
	case "IcmpOutEchos":
		s = cfg.Render(n.IcmpOutEchos)
	case "IcmpOutEchoReps":
		s = cfg.Render(n.IcmpOutEchoReps)
	case "IcmpOutTimestamps":
		s = cfg.Render(n.IcmpOutTimestamps)
	case "IcmpOutTimestampReps":
		s = cfg.Render(n.IcmpOutTimestampReps)
	case "IcmpOutAddrMasks":
		s = cfg.Render(n.IcmpOutAddrMasks)
	case "IcmpOutAddrMaskReps":
		s = cfg.Render(n.IcmpOutAddrMaskReps)
	case "IcmpInType3":
		s = cfg.Render(n.IcmpInType3)
	case "IcmpOutType3":
		s = cfg.Render(n.IcmpOutType3)
	case "TcpActiveOpens":
		s = cfg.Render(n.TcpActiveOpens)
	case "TcpPassiveOpens":
		s = cfg.Render(n.TcpPassiveOpens)
	case "TcpAttemptFails":
		s = cfg.Render(n.TcpAttemptFails)
	case "TcpEstabResets":
		s = cfg.Render(n.TcpEstabResets)
	case "TcpCurrEstab":
		s = cfg.Render(n.TcpCurrEstab)
	case "TcpInSegs":
		s = cfg.Render(n.TcpInSegs)
	case "TcpOutSegs":
		s = cfg.Render(n.TcpOutSegs)
	case "TcpRetransSegs":
		s = cfg.Render(n.TcpRetransSegs)
	case "TcpInErrs":
		s = cfg.Render(n.TcpInErrs)
	case "TcpOutRsts":
		s = cfg.Render(n.TcpOutRsts)
	case "TcpInCsumErrors":
		s = cfg.Render(n.TcpInCsumErrors)
	case "UdpInDatagrams":
		s = cfg.Render(n.UdpInDatagrams)
	case "UdpNoPorts":
		s = cfg.Render(n.UdpNoPorts)
	case "UdpInErrors":
		s = cfg.Render(n.UdpInErrors)
	case "UdpOutDatagrams":
		s = cfg.Render(n.UdpOutDatagrams)
	case "UdpRcvbufErrors":
		s = cfg.Render(n.UdpRcvbufErrors)
	case "UdpSndbufErrors":
		s = cfg.Render(n.UdpSndbufErrors)
	case "UdpInCsumErrors":
		s = cfg.Render(n.UdpInCsumErrors)
	case "UdpIgnoredMulti":
		s = cfg.Render(n.UdpIgnoredMulti)
	case "UdpLiteInDatagrams":
		s = cfg.Render(n.UdpLiteInDatagrams)
	case "UdpLiteNoPorts":
		s = cfg.Render(n.UdpLiteNoPorts)
	case "UdpLiteInErrors":
		s = cfg.Render(n.UdpLiteInErrors)
	case "UdpLiteOutDatagrams":
		s = cfg.Render(n.UdpLiteOutDatagrams)
	case "UdpLiteRcvbufErrors":
		s = cfg.Render(n.UdpLiteRcvbufErrors)
	case "UdpLiteSndbufErrors":
		s = cfg.Render(n.UdpLiteSndbufErrors)
	case "UdpLiteInCsumErrors":
		s = cfg.Render(n.UdpLiteInCsumErrors)
	case "UdpLiteIgnoredMulti":
		s = cfg.Render(n.UdpLiteIgnoredMulti)
	case "Ip6InReceives":
		s = cfg.Render(n.Ip6InReceives)
	case "Ip6InHdrErrors":
		s = cfg.Render(n.Ip6InHdrErrors)
	case "Ip6InTooBigErrors":
		s = cfg.Render(n.Ip6InTooBigErrors)
	case "Ip6InNoRoutes":
		s = cfg.Render(n.Ip6InNoRoutes)
	case "Ip6InAddrErrors":
		s = cfg.Render(n.Ip6InAddrErrors)
	case "Ip6InUnknownProtos":
		s = cfg.Render(n.Ip6InUnknownProtos)
	case "Ip6InTruncatedPkts":
		s = cfg.Render(n.Ip6InTruncatedPkts)
	case "Ip6InDiscards":
		s = cfg.Render(n.Ip6InDiscards)
	case "Ip6InDelivers":
		s = cfg.Render(n.Ip6InDelivers)
	case "Ip6OutForwDatagrams":
		s = cfg.Render(n.Ip6OutForwDatagrams)
	case "Ip6OutRequests":
		s = cfg.Render(n.Ip6OutRequests)
	case "Ip6OutDiscards":
		s = cfg.Render(n.Ip6OutDiscards)
	case "Ip6OutNoRoutes":
		s = cfg.Render(n.Ip6OutNoRoutes)
	case "Ip6ReasmTimeout":
		s = cfg.Render(n.Ip6ReasmTimeout)
	case "Ip6ReasmReqds":
		s = cfg.Render(n.Ip6ReasmReqds)
	case "Ip6ReasmOKs":
		s = cfg.Render(n.Ip6ReasmOKs)
	case "Ip6ReasmFails":
		s = cfg.Render(n.Ip6ReasmFails)
	case "Ip6FragOKs":
		s = cfg.Render(n.Ip6FragOKs)
	case "Ip6FragFails":
		s = cfg.Render(n.Ip6FragFails)
	case "Ip6FragCreates":
		s = cfg.Render(n.Ip6FragCreates)
	case "Ip6InMcastPkts":
		s = cfg.Render(n.Ip6InMcastPkts)
	case "Ip6OutMcastPkts":
		s = cfg.Render(n.Ip6OutMcastPkts)
	case "Ip6InOctets":
		s = cfg.Render(n.Ip6InOctets)
	case "Ip6OutOctets":
		s = cfg.Render(n.Ip6OutOctets)
	case "Ip6InMcastOctets":
		s = cfg.Render(n.Ip6InMcastOctets)
	case "Ip6OutMcastOctets":
		s = cfg.Render(n.Ip6OutMcastOctets)
	case "Ip6InBcastOctets":
		s = cfg.Render(n.Ip6InBcastOctets)
	case "Ip6OutBcastOctets":
		s = cfg.Render(n.Ip6OutBcastOctets)
	case "Ip6InNoECTPkts":
		s = cfg.Render(n.Ip6InNoECTPkts)
	case "Ip6InECT1Pkts":
		s = cfg.Render(n.Ip6InECT1Pkts)
	case "Ip6InECT0Pkts":
		s = cfg.Render(n.Ip6InECT0Pkts)
	case "Ip6InCEPkts":
		s = cfg.Render(n.Ip6InCEPkts)
	case "Icmp6InMsgs":
		s = cfg.Render(n.Icmp6InMsgs)
	case "Icmp6InErrors":
		s = cfg.Render(n.Icmp6InErrors)
	case "Icmp6OutMsgs":
		s = cfg.Render(n.Icmp6OutMsgs)
	case "Icmp6OutErrors":
		s = cfg.Render(n.Icmp6OutErrors)
	case "Icmp6InCsumErrors":
		s = cfg.Render(n.Icmp6InCsumErrors)
	case "Icmp6InDestUnreachs":
		s = cfg.Render(n.Icmp6InDestUnreachs)
	case "Icmp6InPktTooBigs":
		s = cfg.Render(n.Icmp6InPktTooBigs)
	case "Icmp6InTimeExcds":
		s = cfg.Render(n.Icmp6InTimeExcds)
	case "Icmp6InParmProblems":
		s = cfg.Render(n.Icmp6InParmProblems)
	case "Icmp6InEchos":
		s = cfg.Render(n.Icmp6InEchos)
	case "Icmp6InEchoReplies":
		s = cfg.Render(n.Icmp6InEchoReplies)
	case "Icmp6InGroupMembQueries":
		s = cfg.Render(n.Icmp6InGroupMembQueries)
	case "Icmp6InGroupMembResponses":
		s = cfg.Render(n.Icmp6InGroupMembResponses)
	case "Icmp6InGroupMembReductions":
		s = cfg.Render(n.Icmp6InGroupMembReductions)
	case "Icmp6InRouterSolicits":
		s = cfg.Render(n.Icmp6InRouterSolicits)
	case "Icmp6InRouterAdvertisements":
		s = cfg.Render(n.Icmp6InRouterAdvertisements)
	case "Icmp6InNeighborSolicits":
		s = cfg.Render(n.Icmp6InNeighborSolicits)
	case "Icmp6InNeighborAdvertisements":
		s = cfg.Render(n.Icmp6InNeighborAdvertisements)
	case "Icmp6InRedirects":
		s = cfg.Render(n.Icmp6InRedirects)
	case "Icmp6InMLDv2Reports":
		s = cfg.Render(n.Icmp6InMLDv2Reports)
	case "Icmp6OutDestUnreachs":
		s = cfg.Render(n.Icmp6OutDestUnreachs)
	case "Icmp6OutPktTooBigs":
		s = cfg.Render(n.Icmp6OutPktTooBigs)
	case "Icmp6OutTimeExcds":
		s = cfg.Render(n.Icmp6OutTimeExcds)
	case "Icmp6OutParmProblems":
		s = cfg.Render(n.Icmp6OutParmProblems)
	case "Icmp6OutEchos":
		s = cfg.Render(n.Icmp6OutEchos)
	case "Icmp6OutEchoReplies":
		s = cfg.Render(n.Icmp6OutEchoReplies)
	case "Icmp6OutGroupMembQueries":
		s = cfg.Render(n.Icmp6OutGroupMembQueries)
	case "Icmp6OutGroupMembResponses":
		s = cfg.Render(n.Icmp6OutGroupMembResponses)
	case "Icmp6OutGroupMembReductions":
		s = cfg.Render(n.Icmp6OutGroupMembReductions)
	case "Icmp6OutRouterSolicits":
		s = cfg.Render(n.Icmp6OutRouterSolicits)
	case "Icmp6OutRouterAdvertisements":
		s = cfg.Render(n.Icmp6OutRouterAdvertisements)
	case "Icmp6OutNeighborSolicits":
		s = cfg.Render(n.Icmp6OutNeighborSolicits)
	case "Icmp6OutNeighborAdvertisements":
		s = cfg.Render(n.Icmp6OutNeighborAdvertisements)
	case "Icmp6OutRedirects":
		s = cfg.Render(n.Icmp6OutRedirects)
	case "Icmp6OutMLDv2Reports":
		s = cfg.Render(n.Icmp6OutMLDv2Reports)
	case "Icmp6InType1":
		s = cfg.Render(n.Icmp6InType1)
	case "Icmp6InType134":
		s = cfg.Render(n.Icmp6InType134)
	case "Icmp6InType135":
		s = cfg.Render(n.Icmp6InType135)
	case "Icmp6InType136":
		s = cfg.Render(n.Icmp6InType136)
	case "Icmp6InType143":
		s = cfg.Render(n.Icmp6InType143)
	case "Icmp6OutType133":
		s = cfg.Render(n.Icmp6OutType133)
	case "Icmp6OutType135":
		s = cfg.Render(n.Icmp6OutType135)
	case "Icmp6OutType136":
		s = cfg.Render(n.Icmp6OutType136)
	case "Icmp6OutType143":
		s = cfg.Render(n.Icmp6OutType143)
	case "Udp6InDatagrams":
		s = cfg.Render(n.Udp6InDatagrams)
	case "Udp6NoPorts":
		s = cfg.Render(n.Udp6NoPorts)
	case "Udp6InErrors":
		s = cfg.Render(n.Udp6InErrors)
	case "Udp6OutDatagrams":
		s = cfg.Render(n.Udp6OutDatagrams)
	case "Udp6RcvbufErrors":
		s = cfg.Render(n.Udp6RcvbufErrors)
	case "Udp6SndbufErrors":
		s = cfg.Render(n.Udp6SndbufErrors)
	case "Udp6InCsumErrors":
		s = cfg.Render(n.Udp6InCsumErrors)
	case "Udp6IgnoredMulti":
		s = cfg.Render(n.Udp6IgnoredMulti)
	case "UdpLite6InDatagrams":
		s = cfg.Render(n.UdpLite6InDatagrams)
	case "UdpLite6NoPorts":
		s = cfg.Render(n.UdpLite6NoPorts)
	case "UdpLite6InErrors":
		s = cfg.Render(n.UdpLite6InErrors)
	case "UdpLite6OutDatagrams":
		s = cfg.Render(n.UdpLite6OutDatagrams)
	case "UdpLite6RcvbufErrors":
		s = cfg.Render(n.UdpLite6RcvbufErrors)
	case "UdpLite6SndbufErrors":
		s = cfg.Render(n.UdpLite6SndbufErrors)
	case "UdpLite6InCsumErrors":
		s = cfg.Render(n.UdpLite6InCsumErrors)
	case "TcpExtSyncookiesSent":
		s = cfg.Render(n.TcpExtSyncookiesSent)
	case "TcpExtSyncookiesRecv":
		s = cfg.Render(n.TcpExtSyncookiesRecv)
	case "TcpExtSyncookiesFailed":
		s = cfg.Render(n.TcpExtSyncookiesFailed)
	case "TcpExtEmbryonicRsts":
		s = cfg.Render(n.TcpExtEmbryonicRsts)
	case "TcpExtPruneCalled":
		s = cfg.Render(n.TcpExtPruneCalled)
	case "TcpExtRcvPruned":
		s = cfg.Render(n.TcpExtRcvPruned)
	case "TcpExtOfoPruned":
		s = cfg.Render(n.TcpExtOfoPruned)
	case "TcpExtOutOfWindowIcmps":
		s = cfg.Render(n.TcpExtOutOfWindowIcmps)
	case "TcpExtLockDroppedIcmps":
		s = cfg.Render(n.TcpExtLockDroppedIcmps)
	case "TcpExtArpFilter":
		s = cfg.Render(n.TcpExtArpFilter)
	case "TcpExtTW":
		s = cfg.Render(n.TcpExtTW)
	case "TcpExtTWRecycled":
		s = cfg.Render(n.TcpExtTWRecycled)
	case "TcpExtTWKilled":
		s = cfg.Render(n.TcpExtTWKilled)
	case "TcpExtPAWSActive":
		s = cfg.Render(n.TcpExtPAWSActive)
	case "TcpExtPAWSEstab":
		s = cfg.Render(n.TcpExtPAWSEstab)
	case "TcpExtDelayedACKs":
		s = cfg.Render(n.TcpExtDelayedACKs)
	case "TcpExtDelayedACKLocked":
		s = cfg.Render(n.TcpExtDelayedACKLocked)
	case "TcpExtDelayedACKLost":
		s = cfg.Render(n.TcpExtDelayedACKLost)
	case "TcpExtListenOverflows":
		s = cfg.Render(n.TcpExtListenOverflows)
	case "TcpExtListenDrops":
		s = cfg.Render(n.TcpExtListenDrops)
	case "TcpExtTCPHPHits":
		s = cfg.Render(n.TcpExtTCPHPHits)
	case "TcpExtTCPPureAcks":
		s = cfg.Render(n.TcpExtTCPPureAcks)
	case "TcpExtTCPHPAcks":
		s = cfg.Render(n.TcpExtTCPHPAcks)
	case "TcpExtTCPRenoRecovery":
		s = cfg.Render(n.TcpExtTCPRenoRecovery)
	case "TcpExtTCPSackRecovery":
		s = cfg.Render(n.TcpExtTCPSackRecovery)
	case "TcpExtTCPSACKReneging":
		s = cfg.Render(n.TcpExtTCPSACKReneging)
	case "TcpExtTCPSACKReorder":
		s = cfg.Render(n.TcpExtTCPSACKReorder)
	case "TcpExtTCPRenoReorder":
		s = cfg.Render(n.TcpExtTCPRenoReorder)
	case "TcpExtTCPTSReorder":
		s = cfg.Render(n.TcpExtTCPTSReorder)
	case "TcpExtTCPFullUndo":
		s = cfg.Render(n.TcpExtTCPFullUndo)
	case "TcpExtTCPPartialUndo":
		s = cfg.Render(n.TcpExtTCPPartialUndo)
	case "TcpExtTCPDSACKUndo":
		s = cfg.Render(n.TcpExtTCPDSACKUndo)
	case "TcpExtTCPLossUndo":
		s = cfg.Render(n.TcpExtTCPLossUndo)
	case "TcpExtTCPLostRetransmit":
		s = cfg.Render(n.TcpExtTCPLostRetransmit)
	case "TcpExtTCPRenoFailures":
		s = cfg.Render(n.TcpExtTCPRenoFailures)
	case "TcpExtTCPSackFailures":
		s = cfg.Render(n.TcpExtTCPSackFailures)
	case "TcpExtTCPLossFailures":
		s = cfg.Render(n.TcpExtTCPLossFailures)
	case "TcpExtTCPFastRetrans":
		s = cfg.Render(n.TcpExtTCPFastRetrans)
	case "TcpExtTCPSlowStartRetrans":
		s = cfg.Render(n.TcpExtTCPSlowStartRetrans)
	case "TcpExtTCPTimeouts":
		s = cfg.Render(n.TcpExtTCPTimeouts)
	case "TcpExtTCPLossProbes":
		s = cfg.Render(n.TcpExtTCPLossProbes)
	case "TcpExtTCPLossProbeRecovery":
		s = cfg.Render(n.TcpExtTCPLossProbeRecovery)
	case "TcpExtTCPRenoRecoveryFail":
		s = cfg.Render(n.TcpExtTCPRenoRecoveryFail)
	case "TcpExtTCPSackRecoveryFail":
		s = cfg.Render(n.TcpExtTCPSackRecoveryFail)
	case "TcpExtTCPRcvCollapsed":
		s = cfg.Render(n.TcpExtTCPRcvCollapsed)
	case "TcpExtTCPDSACKOldSent":
		s = cfg.Render(n.TcpExtTCPDSACKOldSent)
	case "TcpExtTCPDSACKOfoSent":
		s = cfg.Render(n.TcpExtTCPDSACKOfoSent)
	case "TcpExtTCPDSACKRecv":
		s = cfg.Render(n.TcpExtTCPDSACKRecv)
	case "TcpExtTCPDSACKOfoRecv":
		s = cfg.Render(n.TcpExtTCPDSACKOfoRecv)
	case "TcpExtTCPAbortOnData":
		s = cfg.Render(n.TcpExtTCPAbortOnData)
	case "TcpExtTCPAbortOnClose":
		s = cfg.Render(n.TcpExtTCPAbortOnClose)
	case "TcpExtTCPAbortOnMemory":
		s = cfg.Render(n.TcpExtTCPAbortOnMemory)
	case "TcpExtTCPAbortOnTimeout":
		s = cfg.Render(n.TcpExtTCPAbortOnTimeout)
	case "TcpExtTCPAbortOnLinger":
		s = cfg.Render(n.TcpExtTCPAbortOnLinger)
	case "TcpExtTCPAbortFailed":
		s = cfg.Render(n.TcpExtTCPAbortFailed)
	case "TcpExtTCPMemoryPressures":
		s = cfg.Render(n.TcpExtTCPMemoryPressures)
	case "TcpExtTCPMemoryPressuresChrono":
		s = cfg.Render(n.TcpExtTCPMemoryPressuresChrono)
	case "TcpExtTCPSACKDiscard":
		s = cfg.Render(n.TcpExtTCPSACKDiscard)
	case "TcpExtTCPDSACKIgnoredOld":
		s = cfg.Render(n.TcpExtTCPDSACKIgnoredOld)
	case "TcpExtTCPDSACKIgnoredNoUndo":
		s = cfg.Render(n.TcpExtTCPDSACKIgnoredNoUndo)
	case "TcpExtTCPSpuriousRTOs":
		s = cfg.Render(n.TcpExtTCPSpuriousRTOs)
	case "TcpExtTCPMD5NotFound":
		s = cfg.Render(n.TcpExtTCPMD5NotFound)
	case "TcpExtTCPMD5Unexpected":
		s = cfg.Render(n.TcpExtTCPMD5Unexpected)
	case "TcpExtTCPMD5Failure":
		s = cfg.Render(n.TcpExtTCPMD5Failure)
	case "TcpExtTCPSackShifted":
		s = cfg.Render(n.TcpExtTCPSackShifted)
	case "TcpExtTCPSackMerged":
		s = cfg.Render(n.TcpExtTCPSackMerged)
	case "TcpExtTCPSackShiftFallback":
		s = cfg.Render(n.TcpExtTCPSackShiftFallback)
	case "TcpExtTCPBacklogDrop":
		s = cfg.Render(n.TcpExtTCPBacklogDrop)
	case "TcpExtPFMemallocDrop":
		s = cfg.Render(n.TcpExtPFMemallocDrop)
	case "TcpExtTCPMinTTLDrop":
		s = cfg.Render(n.TcpExtTCPMinTTLDrop)
	case "TcpExtTCPDeferAcceptDrop":
		s = cfg.Render(n.TcpExtTCPDeferAcceptDrop)
	case "TcpExtIPReversePathFilter":
		s = cfg.Render(n.TcpExtIPReversePathFilter)
	case "TcpExtTCPTimeWaitOverflow":
		s = cfg.Render(n.TcpExtTCPTimeWaitOverflow)
	case "TcpExtTCPReqQFullDoCookies":
		s = cfg.Render(n.TcpExtTCPReqQFullDoCookies)
	case "TcpExtTCPReqQFullDrop":
		s = cfg.Render(n.TcpExtTCPReqQFullDrop)
	case "TcpExtTCPRetransFail":
		s = cfg.Render(n.TcpExtTCPRetransFail)
	case "TcpExtTCPRcvCoalesce":
		s = cfg.Render(n.TcpExtTCPRcvCoalesce)
	case "TcpExtTCPRcvQDrop":
		s = cfg.Render(n.TcpExtTCPRcvQDrop)
	case "TcpExtTCPOFOQueue":
		s = cfg.Render(n.TcpExtTCPOFOQueue)
	case "TcpExtTCPOFODrop":
		s = cfg.Render(n.TcpExtTCPOFODrop)
	case "TcpExtTCPOFOMerge":
		s = cfg.Render(n.TcpExtTCPOFOMerge)
	case "TcpExtTCPChallengeACK":
		s = cfg.Render(n.TcpExtTCPChallengeACK)
	case "TcpExtTCPSYNChallenge":
		s = cfg.Render(n.TcpExtTCPSYNChallenge)
	case "TcpExtTCPFastOpenActive":
		s = cfg.Render(n.TcpExtTCPFastOpenActive)
	case "TcpExtTCPFastOpenActiveFail":
		s = cfg.Render(n.TcpExtTCPFastOpenActiveFail)
	case "TcpExtTCPFastOpenPassive":
		s = cfg.Render(n.TcpExtTCPFastOpenPassive)
	case "TcpExtTCPFastOpenPassiveFail":
		s = cfg.Render(n.TcpExtTCPFastOpenPassiveFail)
	case "TcpExtTCPFastOpenListenOverflow":
		s = cfg.Render(n.TcpExtTCPFastOpenListenOverflow)
	case "TcpExtTCPFastOpenCookieReqd":
		s = cfg.Render(n.TcpExtTCPFastOpenCookieReqd)
	case "TcpExtTCPFastOpenBlackhole":
		s = cfg.Render(n.TcpExtTCPFastOpenBlackhole)
	case "TcpExtTCPSpuriousRtxHostQueues":
		s = cfg.Render(n.TcpExtTCPSpuriousRtxHostQueues)
	case "TcpExtBusyPollRxPackets":
		s = cfg.Render(n.TcpExtBusyPollRxPackets)
	case "TcpExtTCPAutoCorking":
		s = cfg.Render(n.TcpExtTCPAutoCorking)
	case "TcpExtTCPFromZeroWindowAdv":
		s = cfg.Render(n.TcpExtTCPFromZeroWindowAdv)
	case "TcpExtTCPToZeroWindowAdv":
		s = cfg.Render(n.TcpExtTCPToZeroWindowAdv)
	case "TcpExtTCPWantZeroWindowAdv":
		s = cfg.Render(n.TcpExtTCPWantZeroWindowAdv)
	case "TcpExtTCPSynRetrans":
		s = cfg.Render(n.TcpExtTCPSynRetrans)
	case "TcpExtTCPOrigDataSent":
		s = cfg.Render(n.TcpExtTCPOrigDataSent)
	case "TcpExtTCPHystartTrainDetect":
		s = cfg.Render(n.TcpExtTCPHystartTrainDetect)
	case "TcpExtTCPHystartTrainCwnd":
		s = cfg.Render(n.TcpExtTCPHystartTrainCwnd)
	case "TcpExtTCPHystartDelayDetect":
		s = cfg.Render(n.TcpExtTCPHystartDelayDetect)
	case "TcpExtTCPHystartDelayCwnd":
		s = cfg.Render(n.TcpExtTCPHystartDelayCwnd)
	case "TcpExtTCPACKSkippedSynRecv":
		s = cfg.Render(n.TcpExtTCPACKSkippedSynRecv)
	case "TcpExtTCPACKSkippedPAWS":
		s = cfg.Render(n.TcpExtTCPACKSkippedPAWS)
	case "TcpExtTCPACKSkippedSeq":
		s = cfg.Render(n.TcpExtTCPACKSkippedSeq)
	case "TcpExtTCPACKSkippedFinWait2":
		s = cfg.Render(n.TcpExtTCPACKSkippedFinWait2)
	case "TcpExtTCPACKSkippedTimeWait":
		s = cfg.Render(n.TcpExtTCPACKSkippedTimeWait)
	case "TcpExtTCPACKSkippedChallenge":
		s = cfg.Render(n.TcpExtTCPACKSkippedChallenge)
	case "TcpExtTCPWinProbe":
		s = cfg.Render(n.TcpExtTCPWinProbe)
	case "TcpExtTCPKeepAlive":
		s = cfg.Render(n.TcpExtTCPKeepAlive)
	case "TcpExtTCPMTUPFail":
		s = cfg.Render(n.TcpExtTCPMTUPFail)
	case "TcpExtTCPMTUPSuccess":
		s = cfg.Render(n.TcpExtTCPMTUPSuccess)
	case "TcpExtTCPWqueueTooBig":
		s = cfg.Render(n.TcpExtTCPWqueueTooBig)
	case "IpExtInNoRoutes":
		s = cfg.Render(n.IpExtInNoRoutes)
	case "IpExtInTruncatedPkts":
		s = cfg.Render(n.IpExtInTruncatedPkts)
	case "IpExtInMcastPkts":
		s = cfg.Render(n.IpExtInMcastPkts)
	case "IpExtOutMcastPkts":
		s = cfg.Render(n.IpExtOutMcastPkts)
	case "IpExtInBcastPkts":
		s = cfg.Render(n.IpExtInBcastPkts)
	case "IpExtOutBcastPkts":
		s = cfg.Render(n.IpExtOutBcastPkts)
	case "IpExtInOctets":
		s = cfg.Render(n.IpExtInOctets)
	case "IpExtOutOctets":
		s = cfg.Render(n.IpExtOutOctets)
	case "IpExtInMcastOctets":
		s = cfg.Render(n.IpExtInMcastOctets)
	case "IpExtOutMcastOctets":
		s = cfg.Render(n.IpExtOutMcastOctets)
	case "IpExtInBcastOctets":
		s = cfg.Render(n.IpExtInBcastOctets)
	case "IpExtOutBcastOctets":
		s = cfg.Render(n.IpExtOutBcastOctets)
	case "IpExtInCsumErrors":
		s = cfg.Render(n.IpExtInCsumErrors)
	case "IpExtInNoECTPkts":
		s = cfg.Render(n.IpExtInNoECTPkts)
	case "IpExtInECT1Pkts":
		s = cfg.Render(n.IpExtInECT1Pkts)
	case "IpExtInECT0Pkts":
		s = cfg.Render(n.IpExtInECT0Pkts)
	case "IpExtInCEPkts":
		s = cfg.Render(n.IpExtInCEPkts)
	case "IpExtReasmOverlaps":
		s = cfg.Render(n.IpExtReasmOverlaps)
	default:
		s = "no " + field + " for netstat stat"
	}
	return s
}

func (netStat *NetStat) Collect(prev, curr *store.Sample) {

	netStat.IpInReceives = getValueOrDefault(curr.Ip.InReceives) - getValueOrDefault(prev.Ip.InReceives)
	netStat.IpInHdrErrors = getValueOrDefault(curr.Ip.InHdrErrors) - getValueOrDefault(prev.Ip.InHdrErrors)
	netStat.IpInAddrErrors = getValueOrDefault(curr.Ip.InAddrErrors) - getValueOrDefault(prev.Ip.InAddrErrors)
	netStat.IpForwDatagrams = getValueOrDefault(curr.Ip.ForwDatagrams) - getValueOrDefault(prev.Ip.ForwDatagrams)
	netStat.IpInUnknownProtos = getValueOrDefault(curr.Ip.InUnknownProtos) - getValueOrDefault(prev.Ip.InUnknownProtos)
	netStat.IpInDiscards = getValueOrDefault(curr.Ip.InDiscards) - getValueOrDefault(prev.Ip.InDiscards)
	netStat.IpInDelivers = getValueOrDefault(curr.Ip.InDelivers) - getValueOrDefault(prev.Ip.InDelivers)
	netStat.IpOutRequests = getValueOrDefault(curr.Ip.OutRequests) - getValueOrDefault(prev.Ip.OutRequests)
	netStat.IpOutDiscards = getValueOrDefault(curr.Ip.OutDiscards) - getValueOrDefault(prev.Ip.OutDiscards)
	netStat.IpOutNoRoutes = getValueOrDefault(curr.Ip.OutNoRoutes) - getValueOrDefault(prev.Ip.OutNoRoutes)
	netStat.IpReasmTimeout = getValueOrDefault(curr.Ip.ReasmTimeout) - getValueOrDefault(prev.Ip.ReasmTimeout)
	netStat.IpReasmReqds = getValueOrDefault(curr.Ip.ReasmReqds) - getValueOrDefault(prev.Ip.ReasmReqds)
	netStat.IpReasmOKs = getValueOrDefault(curr.Ip.ReasmOKs) - getValueOrDefault(prev.Ip.ReasmOKs)
	netStat.IpReasmFails = getValueOrDefault(curr.Ip.ReasmFails) - getValueOrDefault(prev.Ip.ReasmFails)
	netStat.IpFragOKs = getValueOrDefault(curr.Ip.FragOKs) - getValueOrDefault(prev.Ip.FragOKs)
	netStat.IpFragFails = getValueOrDefault(curr.Ip.FragFails) - getValueOrDefault(prev.Ip.FragFails)
	netStat.IpFragCreates = getValueOrDefault(curr.Ip.FragCreates) - getValueOrDefault(prev.Ip.FragCreates)

	netStat.IcmpInMsgs = getValueOrDefault(curr.Icmp.InMsgs) - getValueOrDefault(prev.Icmp.InMsgs)
	netStat.IcmpInErrors = getValueOrDefault(curr.Icmp.InErrors) - getValueOrDefault(prev.Icmp.InErrors)
	netStat.IcmpInCsumErrors = getValueOrDefault(curr.Icmp.InCsumErrors) - getValueOrDefault(prev.Icmp.InCsumErrors)
	netStat.IcmpInDestUnreachs = getValueOrDefault(curr.Icmp.InDestUnreachs) - getValueOrDefault(prev.Icmp.InDestUnreachs)
	netStat.IcmpInTimeExcds = getValueOrDefault(curr.Icmp.InTimeExcds) - getValueOrDefault(prev.Icmp.InTimeExcds)
	netStat.IcmpInParmProbs = getValueOrDefault(curr.Icmp.InParmProbs) - getValueOrDefault(prev.Icmp.InParmProbs)
	netStat.IcmpInSrcQuenchs = getValueOrDefault(curr.Icmp.InSrcQuenchs) - getValueOrDefault(prev.Icmp.InSrcQuenchs)
	netStat.IcmpInRedirects = getValueOrDefault(curr.Icmp.InRedirects) - getValueOrDefault(prev.Icmp.InRedirects)
	netStat.IcmpInEchos = getValueOrDefault(curr.Icmp.InEchos) - getValueOrDefault(prev.Icmp.InEchos)
	netStat.IcmpInEchoReps = getValueOrDefault(curr.Icmp.InEchoReps) - getValueOrDefault(prev.Icmp.InEchoReps)
	netStat.IcmpInTimestamps = getValueOrDefault(curr.Icmp.InTimestamps) - getValueOrDefault(prev.Icmp.InTimestamps)
	netStat.IcmpInTimestampReps = getValueOrDefault(curr.Icmp.InTimestampReps) - getValueOrDefault(prev.Icmp.InTimestampReps)
	netStat.IcmpInAddrMasks = getValueOrDefault(curr.Icmp.InAddrMasks) - getValueOrDefault(prev.Icmp.InAddrMasks)
	netStat.IcmpInAddrMaskReps = getValueOrDefault(curr.Icmp.InAddrMaskReps) - getValueOrDefault(prev.Icmp.InAddrMaskReps)
	netStat.IcmpOutMsgs = getValueOrDefault(curr.Icmp.OutMsgs) - getValueOrDefault(prev.Icmp.OutMsgs)
	netStat.IcmpOutErrors = getValueOrDefault(curr.Icmp.OutErrors) - getValueOrDefault(prev.Icmp.OutErrors)
	netStat.IcmpOutDestUnreachs = getValueOrDefault(curr.Icmp.OutDestUnreachs) - getValueOrDefault(prev.Icmp.OutDestUnreachs)
	netStat.IcmpOutTimeExcds = getValueOrDefault(curr.Icmp.OutTimeExcds) - getValueOrDefault(prev.Icmp.OutTimeExcds)
	netStat.IcmpOutParmProbs = getValueOrDefault(curr.Icmp.OutParmProbs) - getValueOrDefault(prev.Icmp.OutParmProbs)
	netStat.IcmpOutSrcQuenchs = getValueOrDefault(curr.Icmp.OutSrcQuenchs) - getValueOrDefault(prev.Icmp.OutSrcQuenchs)
	netStat.IcmpOutRedirects = getValueOrDefault(curr.Icmp.OutRedirects) - getValueOrDefault(prev.Icmp.OutRedirects)
	netStat.IcmpOutEchos = getValueOrDefault(curr.Icmp.OutEchos) - getValueOrDefault(prev.Icmp.OutEchos)
	netStat.IcmpOutEchoReps = getValueOrDefault(curr.Icmp.OutEchoReps) - getValueOrDefault(prev.Icmp.OutEchoReps)
	netStat.IcmpOutTimestamps = getValueOrDefault(curr.Icmp.OutTimestamps) - getValueOrDefault(prev.Icmp.OutTimestamps)
	netStat.IcmpOutTimestampReps = getValueOrDefault(curr.Icmp.OutTimestampReps) - getValueOrDefault(prev.Icmp.OutTimestampReps)
	netStat.IcmpOutAddrMasks = getValueOrDefault(curr.Icmp.OutAddrMasks) - getValueOrDefault(prev.Icmp.OutAddrMasks)
	netStat.IcmpOutAddrMaskReps = getValueOrDefault(curr.Icmp.OutAddrMaskReps) - getValueOrDefault(prev.Icmp.OutAddrMaskReps)
	netStat.IcmpInType3 = getValueOrDefault(curr.IcmpMsg.InType3) - getValueOrDefault(prev.IcmpMsg.InType3)
	netStat.IcmpOutType3 = getValueOrDefault(curr.IcmpMsg.OutType3) - getValueOrDefault(prev.IcmpMsg.OutType3)

	netStat.TcpActiveOpens = getValueOrDefault(curr.Tcp.ActiveOpens) - getValueOrDefault(prev.Tcp.ActiveOpens)
	netStat.TcpPassiveOpens = getValueOrDefault(curr.Tcp.PassiveOpens) - getValueOrDefault(prev.Tcp.PassiveOpens)
	netStat.TcpAttemptFails = getValueOrDefault(curr.Tcp.AttemptFails) - getValueOrDefault(prev.Tcp.AttemptFails)
	netStat.TcpEstabResets = getValueOrDefault(curr.Tcp.EstabResets) - getValueOrDefault(prev.Tcp.EstabResets)
	netStat.TcpCurrEstab = getValueOrDefault(curr.Tcp.CurrEstab)
	netStat.TcpInSegs = getValueOrDefault(curr.Tcp.InSegs) - getValueOrDefault(prev.Tcp.InSegs)
	netStat.TcpOutSegs = getValueOrDefault(curr.Tcp.OutSegs) - getValueOrDefault(prev.Tcp.OutSegs)
	netStat.TcpRetransSegs = getValueOrDefault(curr.Tcp.RetransSegs) - getValueOrDefault(prev.Tcp.RetransSegs)
	netStat.TcpInErrs = getValueOrDefault(curr.Tcp.InErrs) - getValueOrDefault(prev.Tcp.InErrs)
	netStat.TcpOutRsts = getValueOrDefault(curr.Tcp.OutRsts) - getValueOrDefault(prev.Tcp.OutRsts)
	netStat.TcpInCsumErrors = getValueOrDefault(curr.Tcp.InCsumErrors) - getValueOrDefault(prev.Tcp.InCsumErrors)

	netStat.UdpInDatagrams = getValueOrDefault(curr.Udp.InDatagrams) - getValueOrDefault(prev.Udp.InDatagrams)
	netStat.UdpNoPorts = getValueOrDefault(curr.Udp.NoPorts) - getValueOrDefault(prev.Udp.NoPorts)
	netStat.UdpInErrors = getValueOrDefault(curr.Udp.InErrors) - getValueOrDefault(prev.Udp.InErrors)
	netStat.UdpOutDatagrams = getValueOrDefault(curr.Udp.OutDatagrams) - getValueOrDefault(prev.Udp.OutDatagrams)
	netStat.UdpRcvbufErrors = getValueOrDefault(curr.Udp.RcvbufErrors) - getValueOrDefault(prev.Udp.RcvbufErrors)
	netStat.UdpSndbufErrors = getValueOrDefault(curr.Udp.SndbufErrors) - getValueOrDefault(prev.Udp.SndbufErrors)
	netStat.UdpInCsumErrors = getValueOrDefault(curr.Udp.InCsumErrors) - getValueOrDefault(prev.Udp.InCsumErrors)
	netStat.UdpIgnoredMulti = getValueOrDefault(curr.Udp.IgnoredMulti) - getValueOrDefault(prev.Udp.IgnoredMulti)

	netStat.UdpLiteInDatagrams = getValueOrDefault(curr.UdpLite.InDatagrams) - getValueOrDefault(prev.UdpLite.InDatagrams)
	netStat.UdpLiteNoPorts = getValueOrDefault(curr.UdpLite.NoPorts) - getValueOrDefault(prev.UdpLite.NoPorts)
	netStat.UdpLiteInErrors = getValueOrDefault(curr.UdpLite.InErrors) - getValueOrDefault(prev.UdpLite.InErrors)
	netStat.UdpLiteOutDatagrams = getValueOrDefault(curr.UdpLite.OutDatagrams) - getValueOrDefault(prev.UdpLite.OutDatagrams)
	netStat.UdpLiteRcvbufErrors = getValueOrDefault(curr.UdpLite.RcvbufErrors) - getValueOrDefault(prev.UdpLite.RcvbufErrors)
	netStat.UdpLiteSndbufErrors = getValueOrDefault(curr.UdpLite.SndbufErrors) - getValueOrDefault(prev.UdpLite.SndbufErrors)
	netStat.UdpLiteInCsumErrors = getValueOrDefault(curr.UdpLite.InCsumErrors) - getValueOrDefault(prev.UdpLite.InCsumErrors)
	netStat.UdpLiteIgnoredMulti = getValueOrDefault(curr.UdpLite.IgnoredMulti) - getValueOrDefault(prev.UdpLite.IgnoredMulti)

	netStat.Ip6InReceives = getValueOrDefault(curr.Ip6.InReceives) - getValueOrDefault(prev.Ip6.InReceives)
	netStat.Ip6InHdrErrors = getValueOrDefault(curr.Ip6.InHdrErrors) - getValueOrDefault(prev.Ip6.InHdrErrors)
	netStat.Ip6InTooBigErrors = getValueOrDefault(curr.Ip6.InTooBigErrors) - getValueOrDefault(prev.Ip6.InTooBigErrors)
	netStat.Ip6InNoRoutes = getValueOrDefault(curr.Ip6.InNoRoutes) - getValueOrDefault(prev.Ip6.InNoRoutes)
	netStat.Ip6InAddrErrors = getValueOrDefault(curr.Ip6.InAddrErrors) - getValueOrDefault(prev.Ip6.InAddrErrors)
	netStat.Ip6InUnknownProtos = getValueOrDefault(curr.Ip6.InUnknownProtos) - getValueOrDefault(prev.Ip6.InUnknownProtos)
	netStat.Ip6InTruncatedPkts = getValueOrDefault(curr.Ip6.InTruncatedPkts) - getValueOrDefault(prev.Ip6.InTruncatedPkts)
	netStat.Ip6InDiscards = getValueOrDefault(curr.Ip6.InDiscards) - getValueOrDefault(prev.Ip6.InDiscards)
	netStat.Ip6InDelivers = getValueOrDefault(curr.Ip6.InDelivers) - getValueOrDefault(prev.Ip6.InDelivers)
	netStat.Ip6OutForwDatagrams = getValueOrDefault(curr.Ip6.OutForwDatagrams) - getValueOrDefault(prev.Ip6.OutForwDatagrams)
	netStat.Ip6OutRequests = getValueOrDefault(curr.Ip6.OutRequests) - getValueOrDefault(prev.Ip6.OutRequests)
	netStat.Ip6OutDiscards = getValueOrDefault(curr.Ip6.OutDiscards) - getValueOrDefault(prev.Ip6.OutDiscards)
	netStat.Ip6OutNoRoutes = getValueOrDefault(curr.Ip6.OutNoRoutes) - getValueOrDefault(prev.Ip6.OutNoRoutes)
	netStat.Ip6ReasmTimeout = getValueOrDefault(curr.Ip6.ReasmTimeout) - getValueOrDefault(prev.Ip6.ReasmTimeout)
	netStat.Ip6ReasmReqds = getValueOrDefault(curr.Ip6.ReasmReqds) - getValueOrDefault(prev.Ip6.ReasmReqds)
	netStat.Ip6ReasmOKs = getValueOrDefault(curr.Ip6.ReasmOKs) - getValueOrDefault(prev.Ip6.ReasmOKs)
	netStat.Ip6ReasmFails = getValueOrDefault(curr.Ip6.ReasmFails) - getValueOrDefault(prev.Ip6.ReasmFails)
	netStat.Ip6FragOKs = getValueOrDefault(curr.Ip6.FragOKs) - getValueOrDefault(prev.Ip6.FragOKs)
	netStat.Ip6FragFails = getValueOrDefault(curr.Ip6.FragFails) - getValueOrDefault(prev.Ip6.FragFails)
	netStat.Ip6FragCreates = getValueOrDefault(curr.Ip6.FragCreates) - getValueOrDefault(prev.Ip6.FragCreates)
	netStat.Ip6InMcastPkts = getValueOrDefault(curr.Ip6.InMcastPkts) - getValueOrDefault(prev.Ip6.InMcastPkts)
	netStat.Ip6OutMcastPkts = getValueOrDefault(curr.Ip6.OutMcastPkts) - getValueOrDefault(prev.Ip6.OutMcastPkts)
	netStat.Ip6InOctets = getValueOrDefault(curr.Ip6.InOctets) - getValueOrDefault(prev.Ip6.InOctets)
	netStat.Ip6OutOctets = getValueOrDefault(curr.Ip6.OutOctets) - getValueOrDefault(prev.Ip6.OutOctets)
	netStat.Ip6InMcastOctets = getValueOrDefault(curr.Ip6.InMcastOctets) - getValueOrDefault(prev.Ip6.InMcastOctets)
	netStat.Ip6OutMcastOctets = getValueOrDefault(curr.Ip6.OutMcastOctets) - getValueOrDefault(prev.Ip6.OutMcastOctets)
	netStat.Ip6InBcastOctets = getValueOrDefault(curr.Ip6.InBcastOctets) - getValueOrDefault(prev.Ip6.InBcastOctets)
	netStat.Ip6OutBcastOctets = getValueOrDefault(curr.Ip6.OutBcastOctets) - getValueOrDefault(prev.Ip6.OutBcastOctets)
	netStat.Ip6InNoECTPkts = getValueOrDefault(curr.Ip6.InNoECTPkts) - getValueOrDefault(prev.Ip6.InNoECTPkts)
	netStat.Ip6InECT1Pkts = getValueOrDefault(curr.Ip6.InECT1Pkts) - getValueOrDefault(prev.Ip6.InECT1Pkts)
	netStat.Ip6InECT0Pkts = getValueOrDefault(curr.Ip6.InECT0Pkts) - getValueOrDefault(prev.Ip6.InECT0Pkts)
	netStat.Ip6InCEPkts = getValueOrDefault(curr.Ip6.InCEPkts) - getValueOrDefault(prev.Ip6.InCEPkts)

	netStat.Icmp6InMsgs = getValueOrDefault(curr.Icmp6.InMsgs) - getValueOrDefault(prev.Icmp6.InMsgs)
	netStat.Icmp6InErrors = getValueOrDefault(curr.Icmp6.InErrors) - getValueOrDefault(prev.Icmp6.InErrors)
	netStat.Icmp6OutMsgs = getValueOrDefault(curr.Icmp6.OutMsgs) - getValueOrDefault(prev.Icmp6.OutMsgs)
	netStat.Icmp6OutErrors = getValueOrDefault(curr.Icmp6.OutErrors) - getValueOrDefault(prev.Icmp6.OutErrors)
	netStat.Icmp6InCsumErrors = getValueOrDefault(curr.Icmp6.InCsumErrors) - getValueOrDefault(prev.Icmp6.InCsumErrors)
	netStat.Icmp6InDestUnreachs = getValueOrDefault(curr.Icmp6.InDestUnreachs) - getValueOrDefault(prev.Icmp6.InDestUnreachs)
	netStat.Icmp6InPktTooBigs = getValueOrDefault(curr.Icmp6.InPktTooBigs) - getValueOrDefault(prev.Icmp6.InPktTooBigs)
	netStat.Icmp6InTimeExcds = getValueOrDefault(curr.Icmp6.InTimeExcds) - getValueOrDefault(prev.Icmp6.InTimeExcds)
	netStat.Icmp6InParmProblems = getValueOrDefault(curr.Icmp6.InParmProblems) - getValueOrDefault(prev.Icmp6.InParmProblems)
	netStat.Icmp6InEchos = getValueOrDefault(curr.Icmp6.InEchos) - getValueOrDefault(prev.Icmp6.InEchos)
	netStat.Icmp6InEchoReplies = getValueOrDefault(curr.Icmp6.InEchoReplies) - getValueOrDefault(prev.Icmp6.InEchoReplies)
	netStat.Icmp6InGroupMembQueries = getValueOrDefault(curr.Icmp6.InGroupMembQueries) - getValueOrDefault(prev.Icmp6.InGroupMembQueries)
	netStat.Icmp6InGroupMembResponses = getValueOrDefault(curr.Icmp6.InGroupMembResponses) - getValueOrDefault(prev.Icmp6.InGroupMembResponses)
	netStat.Icmp6InGroupMembReductions = getValueOrDefault(curr.Icmp6.InGroupMembReductions) - getValueOrDefault(prev.Icmp6.InGroupMembReductions)
	netStat.Icmp6InRouterSolicits = getValueOrDefault(curr.Icmp6.InRouterSolicits) - getValueOrDefault(prev.Icmp6.InRouterSolicits)
	netStat.Icmp6InRouterAdvertisements = getValueOrDefault(curr.Icmp6.InRouterAdvertisements) - getValueOrDefault(prev.Icmp6.InRouterAdvertisements)
	netStat.Icmp6InNeighborSolicits = getValueOrDefault(curr.Icmp6.InNeighborSolicits) - getValueOrDefault(prev.Icmp6.InNeighborSolicits)
	netStat.Icmp6InNeighborAdvertisements = getValueOrDefault(curr.Icmp6.InNeighborAdvertisements) - getValueOrDefault(prev.Icmp6.InNeighborAdvertisements)
	netStat.Icmp6InRedirects = getValueOrDefault(curr.Icmp6.InRedirects) - getValueOrDefault(prev.Icmp6.InRedirects)
	netStat.Icmp6InMLDv2Reports = getValueOrDefault(curr.Icmp6.InMLDv2Reports) - getValueOrDefault(prev.Icmp6.InMLDv2Reports)
	netStat.Icmp6OutDestUnreachs = getValueOrDefault(curr.Icmp6.OutDestUnreachs) - getValueOrDefault(prev.Icmp6.OutDestUnreachs)
	netStat.Icmp6OutPktTooBigs = getValueOrDefault(curr.Icmp6.OutPktTooBigs) - getValueOrDefault(prev.Icmp6.OutPktTooBigs)
	netStat.Icmp6OutTimeExcds = getValueOrDefault(curr.Icmp6.OutTimeExcds) - getValueOrDefault(prev.Icmp6.OutTimeExcds)
	netStat.Icmp6OutParmProblems = getValueOrDefault(curr.Icmp6.OutParmProblems) - getValueOrDefault(prev.Icmp6.OutParmProblems)
	netStat.Icmp6OutEchos = getValueOrDefault(curr.Icmp6.OutEchos) - getValueOrDefault(prev.Icmp6.OutEchos)
	netStat.Icmp6OutEchoReplies = getValueOrDefault(curr.Icmp6.OutEchoReplies) - getValueOrDefault(prev.Icmp6.OutEchoReplies)
	netStat.Icmp6OutGroupMembQueries = getValueOrDefault(curr.Icmp6.OutGroupMembQueries) - getValueOrDefault(prev.Icmp6.OutGroupMembQueries)
	netStat.Icmp6OutGroupMembResponses = getValueOrDefault(curr.Icmp6.OutGroupMembResponses) - getValueOrDefault(prev.Icmp6.OutGroupMembResponses)
	netStat.Icmp6OutGroupMembReductions = getValueOrDefault(curr.Icmp6.OutGroupMembReductions) - getValueOrDefault(prev.Icmp6.OutGroupMembReductions)
	netStat.Icmp6OutRouterSolicits = getValueOrDefault(curr.Icmp6.OutRouterSolicits) - getValueOrDefault(prev.Icmp6.OutRouterSolicits)
	netStat.Icmp6OutRouterAdvertisements = getValueOrDefault(curr.Icmp6.OutRouterAdvertisements) - getValueOrDefault(prev.Icmp6.OutRouterAdvertisements)
	netStat.Icmp6OutNeighborSolicits = getValueOrDefault(curr.Icmp6.OutNeighborSolicits) - getValueOrDefault(prev.Icmp6.OutNeighborSolicits)
	netStat.Icmp6OutNeighborAdvertisements = getValueOrDefault(curr.Icmp6.OutNeighborAdvertisements) - getValueOrDefault(prev.Icmp6.OutNeighborAdvertisements)
	netStat.Icmp6OutRedirects = getValueOrDefault(curr.Icmp6.OutRedirects) - getValueOrDefault(prev.Icmp6.OutRedirects)
	netStat.Icmp6OutMLDv2Reports = getValueOrDefault(curr.Icmp6.OutMLDv2Reports) - getValueOrDefault(prev.Icmp6.OutMLDv2Reports)
	netStat.Icmp6InType1 = getValueOrDefault(curr.Icmp6.InType1) - getValueOrDefault(prev.Icmp6.InType1)
	netStat.Icmp6InType134 = getValueOrDefault(curr.Icmp6.InType134) - getValueOrDefault(prev.Icmp6.InType134)
	netStat.Icmp6InType135 = getValueOrDefault(curr.Icmp6.InType135) - getValueOrDefault(prev.Icmp6.InType135)
	netStat.Icmp6InType136 = getValueOrDefault(curr.Icmp6.InType136) - getValueOrDefault(prev.Icmp6.InType136)
	netStat.Icmp6InType143 = getValueOrDefault(curr.Icmp6.InType143) - getValueOrDefault(prev.Icmp6.InType143)
	netStat.Icmp6OutType133 = getValueOrDefault(curr.Icmp6.OutType133) - getValueOrDefault(prev.Icmp6.OutType133)
	netStat.Icmp6OutType135 = getValueOrDefault(curr.Icmp6.OutType135) - getValueOrDefault(prev.Icmp6.OutType135)
	netStat.Icmp6OutType136 = getValueOrDefault(curr.Icmp6.OutType136) - getValueOrDefault(prev.Icmp6.OutType136)
	netStat.Icmp6OutType143 = getValueOrDefault(curr.Icmp6.OutType143) - getValueOrDefault(prev.Icmp6.OutType143)

	netStat.Udp6InDatagrams = getValueOrDefault(curr.Udp6.InDatagrams) - getValueOrDefault(prev.Udp6.InDatagrams)
	netStat.Udp6NoPorts = getValueOrDefault(curr.Udp6.NoPorts) - getValueOrDefault(prev.Udp6.NoPorts)
	netStat.Udp6InErrors = getValueOrDefault(curr.Udp6.InErrors) - getValueOrDefault(prev.Udp6.InErrors)
	netStat.Udp6OutDatagrams = getValueOrDefault(curr.Udp6.OutDatagrams) - getValueOrDefault(prev.Udp6.OutDatagrams)
	netStat.Udp6RcvbufErrors = getValueOrDefault(curr.Udp6.RcvbufErrors) - getValueOrDefault(prev.Udp6.RcvbufErrors)
	netStat.Udp6SndbufErrors = getValueOrDefault(curr.Udp6.SndbufErrors) - getValueOrDefault(prev.Udp6.SndbufErrors)
	netStat.Udp6InCsumErrors = getValueOrDefault(curr.Udp6.InCsumErrors) - getValueOrDefault(prev.Udp6.InCsumErrors)
	netStat.Udp6IgnoredMulti = getValueOrDefault(curr.Udp6.IgnoredMulti) - getValueOrDefault(prev.Udp6.IgnoredMulti)

	netStat.UdpLite6InDatagrams = getValueOrDefault(curr.UdpLite6.InDatagrams) - getValueOrDefault(prev.UdpLite6.InDatagrams)
	netStat.UdpLite6NoPorts = getValueOrDefault(curr.UdpLite6.NoPorts) - getValueOrDefault(prev.UdpLite6.NoPorts)
	netStat.UdpLite6InErrors = getValueOrDefault(curr.UdpLite6.InErrors) - getValueOrDefault(prev.UdpLite6.InErrors)
	netStat.UdpLite6OutDatagrams = getValueOrDefault(curr.UdpLite6.OutDatagrams) - getValueOrDefault(prev.UdpLite6.OutDatagrams)
	netStat.UdpLite6RcvbufErrors = getValueOrDefault(curr.UdpLite6.RcvbufErrors) - getValueOrDefault(prev.UdpLite6.RcvbufErrors)
	netStat.UdpLite6SndbufErrors = getValueOrDefault(curr.UdpLite6.SndbufErrors) - getValueOrDefault(prev.UdpLite6.SndbufErrors)
	netStat.UdpLite6InCsumErrors = getValueOrDefault(curr.UdpLite6.InCsumErrors) - getValueOrDefault(prev.UdpLite6.InCsumErrors)

	netStat.TcpExtSyncookiesSent = getValueOrDefault(curr.TcpExt.SyncookiesSent) - getValueOrDefault(prev.TcpExt.SyncookiesSent)
	netStat.TcpExtSyncookiesRecv = getValueOrDefault(curr.TcpExt.SyncookiesRecv) - getValueOrDefault(prev.TcpExt.SyncookiesRecv)
	netStat.TcpExtSyncookiesFailed = getValueOrDefault(curr.TcpExt.SyncookiesFailed) - getValueOrDefault(prev.TcpExt.SyncookiesFailed)
	netStat.TcpExtEmbryonicRsts = getValueOrDefault(curr.TcpExt.EmbryonicRsts) - getValueOrDefault(prev.TcpExt.EmbryonicRsts)
	netStat.TcpExtPruneCalled = getValueOrDefault(curr.TcpExt.PruneCalled) - getValueOrDefault(prev.TcpExt.PruneCalled)
	netStat.TcpExtRcvPruned = getValueOrDefault(curr.TcpExt.RcvPruned) - getValueOrDefault(prev.TcpExt.RcvPruned)
	netStat.TcpExtOfoPruned = getValueOrDefault(curr.TcpExt.OfoPruned) - getValueOrDefault(prev.TcpExt.OfoPruned)
	netStat.TcpExtOutOfWindowIcmps = getValueOrDefault(curr.TcpExt.OutOfWindowIcmps) - getValueOrDefault(prev.TcpExt.OutOfWindowIcmps)
	netStat.TcpExtLockDroppedIcmps = getValueOrDefault(curr.TcpExt.LockDroppedIcmps) - getValueOrDefault(prev.TcpExt.LockDroppedIcmps)
	netStat.TcpExtArpFilter = getValueOrDefault(curr.TcpExt.ArpFilter) - getValueOrDefault(prev.TcpExt.ArpFilter)
	netStat.TcpExtTW = getValueOrDefault(curr.TcpExt.TW) - getValueOrDefault(prev.TcpExt.TW)
	netStat.TcpExtTWRecycled = getValueOrDefault(curr.TcpExt.TWRecycled) - getValueOrDefault(prev.TcpExt.TWRecycled)
	netStat.TcpExtTWKilled = getValueOrDefault(curr.TcpExt.TWKilled) - getValueOrDefault(prev.TcpExt.TWKilled)
	netStat.TcpExtPAWSActive = getValueOrDefault(curr.TcpExt.PAWSActive) - getValueOrDefault(prev.TcpExt.PAWSActive)
	netStat.TcpExtPAWSEstab = getValueOrDefault(curr.TcpExt.PAWSEstab) - getValueOrDefault(prev.TcpExt.PAWSEstab)
	netStat.TcpExtDelayedACKs = getValueOrDefault(curr.TcpExt.DelayedACKs) - getValueOrDefault(prev.TcpExt.DelayedACKs)
	netStat.TcpExtDelayedACKLocked = getValueOrDefault(curr.TcpExt.DelayedACKLocked) - getValueOrDefault(prev.TcpExt.DelayedACKLocked)
	netStat.TcpExtDelayedACKLost = getValueOrDefault(curr.TcpExt.DelayedACKLost) - getValueOrDefault(prev.TcpExt.DelayedACKLost)
	netStat.TcpExtListenOverflows = getValueOrDefault(curr.TcpExt.ListenOverflows) - getValueOrDefault(prev.TcpExt.ListenOverflows)
	netStat.TcpExtListenDrops = getValueOrDefault(curr.TcpExt.ListenDrops) - getValueOrDefault(prev.TcpExt.ListenDrops)
	netStat.TcpExtTCPHPHits = getValueOrDefault(curr.TcpExt.TCPHPHits) - getValueOrDefault(prev.TcpExt.TCPHPHits)
	netStat.TcpExtTCPPureAcks = getValueOrDefault(curr.TcpExt.TCPPureAcks) - getValueOrDefault(prev.TcpExt.TCPPureAcks)
	netStat.TcpExtTCPHPAcks = getValueOrDefault(curr.TcpExt.TCPHPAcks) - getValueOrDefault(prev.TcpExt.TCPHPAcks)
	netStat.TcpExtTCPRenoRecovery = getValueOrDefault(curr.TcpExt.TCPRenoRecovery) - getValueOrDefault(prev.TcpExt.TCPRenoRecovery)
	netStat.TcpExtTCPSackRecovery = getValueOrDefault(curr.TcpExt.TCPSackRecovery) - getValueOrDefault(prev.TcpExt.TCPSackRecovery)
	netStat.TcpExtTCPSACKReneging = getValueOrDefault(curr.TcpExt.TCPSACKReneging) - getValueOrDefault(prev.TcpExt.TCPSACKReneging)
	netStat.TcpExtTCPSACKReorder = getValueOrDefault(curr.TcpExt.TCPSACKReorder) - getValueOrDefault(prev.TcpExt.TCPSACKReorder)
	netStat.TcpExtTCPRenoReorder = getValueOrDefault(curr.TcpExt.TCPRenoReorder) - getValueOrDefault(prev.TcpExt.TCPRenoReorder)
	netStat.TcpExtTCPTSReorder = getValueOrDefault(curr.TcpExt.TCPTSReorder) - getValueOrDefault(prev.TcpExt.TCPTSReorder)
	netStat.TcpExtTCPFullUndo = getValueOrDefault(curr.TcpExt.TCPFullUndo) - getValueOrDefault(prev.TcpExt.TCPFullUndo)
	netStat.TcpExtTCPPartialUndo = getValueOrDefault(curr.TcpExt.TCPPartialUndo) - getValueOrDefault(prev.TcpExt.TCPPartialUndo)
	netStat.TcpExtTCPDSACKUndo = getValueOrDefault(curr.TcpExt.TCPDSACKUndo) - getValueOrDefault(prev.TcpExt.TCPDSACKUndo)
	netStat.TcpExtTCPLossUndo = getValueOrDefault(curr.TcpExt.TCPLossUndo) - getValueOrDefault(prev.TcpExt.TCPLossUndo)
	netStat.TcpExtTCPLostRetransmit = getValueOrDefault(curr.TcpExt.TCPLostRetransmit) - getValueOrDefault(prev.TcpExt.TCPLostRetransmit)
	netStat.TcpExtTCPRenoFailures = getValueOrDefault(curr.TcpExt.TCPRenoFailures) - getValueOrDefault(prev.TcpExt.TCPRenoFailures)
	netStat.TcpExtTCPSackFailures = getValueOrDefault(curr.TcpExt.TCPSackFailures) - getValueOrDefault(prev.TcpExt.TCPSackFailures)
	netStat.TcpExtTCPLossFailures = getValueOrDefault(curr.TcpExt.TCPLossFailures) - getValueOrDefault(prev.TcpExt.TCPLossFailures)
	netStat.TcpExtTCPFastRetrans = getValueOrDefault(curr.TcpExt.TCPFastRetrans) - getValueOrDefault(prev.TcpExt.TCPFastRetrans)
	netStat.TcpExtTCPSlowStartRetrans = getValueOrDefault(curr.TcpExt.TCPSlowStartRetrans) - getValueOrDefault(prev.TcpExt.TCPSlowStartRetrans)
	netStat.TcpExtTCPTimeouts = getValueOrDefault(curr.TcpExt.TCPTimeouts) - getValueOrDefault(prev.TcpExt.TCPTimeouts)
	netStat.TcpExtTCPLossProbes = getValueOrDefault(curr.TcpExt.TCPLossProbes) - getValueOrDefault(prev.TcpExt.TCPLossProbes)
	netStat.TcpExtTCPLossProbeRecovery = getValueOrDefault(curr.TcpExt.TCPLossProbeRecovery) - getValueOrDefault(prev.TcpExt.TCPLossProbeRecovery)
	netStat.TcpExtTCPRenoRecoveryFail = getValueOrDefault(curr.TcpExt.TCPRenoRecoveryFail) - getValueOrDefault(prev.TcpExt.TCPRenoRecoveryFail)
	netStat.TcpExtTCPSackRecoveryFail = getValueOrDefault(curr.TcpExt.TCPSackRecoveryFail) - getValueOrDefault(prev.TcpExt.TCPSackRecoveryFail)
	netStat.TcpExtTCPRcvCollapsed = getValueOrDefault(curr.TcpExt.TCPRcvCollapsed) - getValueOrDefault(prev.TcpExt.TCPRcvCollapsed)
	netStat.TcpExtTCPDSACKOldSent = getValueOrDefault(curr.TcpExt.TCPDSACKOldSent) - getValueOrDefault(prev.TcpExt.TCPDSACKOldSent)
	netStat.TcpExtTCPDSACKOfoSent = getValueOrDefault(curr.TcpExt.TCPDSACKOfoSent) - getValueOrDefault(prev.TcpExt.TCPDSACKOfoSent)
	netStat.TcpExtTCPDSACKRecv = getValueOrDefault(curr.TcpExt.TCPDSACKRecv) - getValueOrDefault(prev.TcpExt.TCPDSACKRecv)
	netStat.TcpExtTCPDSACKOfoRecv = getValueOrDefault(curr.TcpExt.TCPDSACKOfoRecv) - getValueOrDefault(prev.TcpExt.TCPDSACKOfoRecv)
	netStat.TcpExtTCPAbortOnData = getValueOrDefault(curr.TcpExt.TCPAbortOnData) - getValueOrDefault(prev.TcpExt.TCPAbortOnData)
	netStat.TcpExtTCPAbortOnClose = getValueOrDefault(curr.TcpExt.TCPAbortOnClose) - getValueOrDefault(prev.TcpExt.TCPAbortOnClose)
	netStat.TcpExtTCPAbortOnMemory = getValueOrDefault(curr.TcpExt.TCPAbortOnMemory) - getValueOrDefault(prev.TcpExt.TCPAbortOnMemory)
	netStat.TcpExtTCPAbortOnTimeout = getValueOrDefault(curr.TcpExt.TCPAbortOnTimeout) - getValueOrDefault(prev.TcpExt.TCPAbortOnTimeout)
	netStat.TcpExtTCPAbortOnLinger = getValueOrDefault(curr.TcpExt.TCPAbortOnLinger) - getValueOrDefault(prev.TcpExt.TCPAbortOnLinger)
	netStat.TcpExtTCPAbortFailed = getValueOrDefault(curr.TcpExt.TCPAbortFailed) - getValueOrDefault(prev.TcpExt.TCPAbortFailed)
	netStat.TcpExtTCPMemoryPressures = getValueOrDefault(curr.TcpExt.TCPMemoryPressures) - getValueOrDefault(prev.TcpExt.TCPMemoryPressures)
	netStat.TcpExtTCPMemoryPressuresChrono = getValueOrDefault(curr.TcpExt.TCPMemoryPressuresChrono) - getValueOrDefault(prev.TcpExt.TCPMemoryPressuresChrono)
	netStat.TcpExtTCPSACKDiscard = getValueOrDefault(curr.TcpExt.TCPSACKDiscard) - getValueOrDefault(prev.TcpExt.TCPSACKDiscard)
	netStat.TcpExtTCPDSACKIgnoredOld = getValueOrDefault(curr.TcpExt.TCPDSACKIgnoredOld) - getValueOrDefault(prev.TcpExt.TCPDSACKIgnoredOld)
	netStat.TcpExtTCPDSACKIgnoredNoUndo = getValueOrDefault(curr.TcpExt.TCPDSACKIgnoredNoUndo) - getValueOrDefault(prev.TcpExt.TCPDSACKIgnoredNoUndo)
	netStat.TcpExtTCPSpuriousRTOs = getValueOrDefault(curr.TcpExt.TCPSpuriousRTOs) - getValueOrDefault(prev.TcpExt.TCPSpuriousRTOs)
	netStat.TcpExtTCPMD5NotFound = getValueOrDefault(curr.TcpExt.TCPMD5NotFound) - getValueOrDefault(prev.TcpExt.TCPMD5NotFound)
	netStat.TcpExtTCPMD5Unexpected = getValueOrDefault(curr.TcpExt.TCPMD5Unexpected) - getValueOrDefault(prev.TcpExt.TCPMD5Unexpected)
	netStat.TcpExtTCPMD5Failure = getValueOrDefault(curr.TcpExt.TCPMD5Failure) - getValueOrDefault(prev.TcpExt.TCPMD5Failure)
	netStat.TcpExtTCPSackShifted = getValueOrDefault(curr.TcpExt.TCPSackShifted) - getValueOrDefault(prev.TcpExt.TCPSackShifted)
	netStat.TcpExtTCPSackMerged = getValueOrDefault(curr.TcpExt.TCPSackMerged) - getValueOrDefault(prev.TcpExt.TCPSackMerged)
	netStat.TcpExtTCPSackShiftFallback = getValueOrDefault(curr.TcpExt.TCPSackShiftFallback) - getValueOrDefault(prev.TcpExt.TCPSackShiftFallback)
	netStat.TcpExtTCPBacklogDrop = getValueOrDefault(curr.TcpExt.TCPBacklogDrop) - getValueOrDefault(prev.TcpExt.TCPBacklogDrop)
	netStat.TcpExtPFMemallocDrop = getValueOrDefault(curr.TcpExt.PFMemallocDrop) - getValueOrDefault(prev.TcpExt.PFMemallocDrop)
	netStat.TcpExtTCPMinTTLDrop = getValueOrDefault(curr.TcpExt.TCPMinTTLDrop) - getValueOrDefault(prev.TcpExt.TCPMinTTLDrop)
	netStat.TcpExtTCPDeferAcceptDrop = getValueOrDefault(curr.TcpExt.TCPDeferAcceptDrop) - getValueOrDefault(prev.TcpExt.TCPDeferAcceptDrop)
	netStat.TcpExtIPReversePathFilter = getValueOrDefault(curr.TcpExt.IPReversePathFilter) - getValueOrDefault(prev.TcpExt.IPReversePathFilter)
	netStat.TcpExtTCPTimeWaitOverflow = getValueOrDefault(curr.TcpExt.TCPTimeWaitOverflow) - getValueOrDefault(prev.TcpExt.TCPTimeWaitOverflow)
	netStat.TcpExtTCPReqQFullDoCookies = getValueOrDefault(curr.TcpExt.TCPReqQFullDoCookies) - getValueOrDefault(prev.TcpExt.TCPReqQFullDoCookies)
	netStat.TcpExtTCPReqQFullDrop = getValueOrDefault(curr.TcpExt.TCPReqQFullDrop) - getValueOrDefault(prev.TcpExt.TCPReqQFullDrop)
	netStat.TcpExtTCPRetransFail = getValueOrDefault(curr.TcpExt.TCPRetransFail) - getValueOrDefault(prev.TcpExt.TCPRetransFail)
	netStat.TcpExtTCPRcvCoalesce = getValueOrDefault(curr.TcpExt.TCPRcvCoalesce) - getValueOrDefault(prev.TcpExt.TCPRcvCoalesce)
	netStat.TcpExtTCPRcvQDrop = getValueOrDefault(curr.TcpExt.TCPRcvQDrop) - getValueOrDefault(prev.TcpExt.TCPRcvQDrop)
	netStat.TcpExtTCPOFOQueue = getValueOrDefault(curr.TcpExt.TCPOFOQueue) - getValueOrDefault(prev.TcpExt.TCPOFOQueue)
	netStat.TcpExtTCPOFODrop = getValueOrDefault(curr.TcpExt.TCPOFODrop) - getValueOrDefault(prev.TcpExt.TCPOFODrop)
	netStat.TcpExtTCPOFOMerge = getValueOrDefault(curr.TcpExt.TCPOFOMerge) - getValueOrDefault(prev.TcpExt.TCPOFOMerge)
	netStat.TcpExtTCPChallengeACK = getValueOrDefault(curr.TcpExt.TCPChallengeACK) - getValueOrDefault(prev.TcpExt.TCPChallengeACK)
	netStat.TcpExtTCPSYNChallenge = getValueOrDefault(curr.TcpExt.TCPSYNChallenge) - getValueOrDefault(prev.TcpExt.TCPSYNChallenge)
	netStat.TcpExtTCPFastOpenActive = getValueOrDefault(curr.TcpExt.TCPFastOpenActive) - getValueOrDefault(prev.TcpExt.TCPFastOpenActive)
	netStat.TcpExtTCPFastOpenActiveFail = getValueOrDefault(curr.TcpExt.TCPFastOpenActiveFail) - getValueOrDefault(prev.TcpExt.TCPFastOpenActiveFail)
	netStat.TcpExtTCPFastOpenPassive = getValueOrDefault(curr.TcpExt.TCPFastOpenPassive) - getValueOrDefault(prev.TcpExt.TCPFastOpenPassive)
	netStat.TcpExtTCPFastOpenPassiveFail = getValueOrDefault(curr.TcpExt.TCPFastOpenPassiveFail) - getValueOrDefault(prev.TcpExt.TCPFastOpenPassiveFail)
	netStat.TcpExtTCPFastOpenListenOverflow = getValueOrDefault(curr.TcpExt.TCPFastOpenListenOverflow) - getValueOrDefault(prev.TcpExt.TCPFastOpenListenOverflow)
	netStat.TcpExtTCPFastOpenCookieReqd = getValueOrDefault(curr.TcpExt.TCPFastOpenCookieReqd) - getValueOrDefault(prev.TcpExt.TCPFastOpenCookieReqd)
	netStat.TcpExtTCPFastOpenBlackhole = getValueOrDefault(curr.TcpExt.TCPFastOpenBlackhole) - getValueOrDefault(prev.TcpExt.TCPFastOpenBlackhole)
	netStat.TcpExtTCPSpuriousRtxHostQueues = getValueOrDefault(curr.TcpExt.TCPSpuriousRtxHostQueues) - getValueOrDefault(prev.TcpExt.TCPSpuriousRtxHostQueues)
	netStat.TcpExtBusyPollRxPackets = getValueOrDefault(curr.TcpExt.BusyPollRxPackets) - getValueOrDefault(prev.TcpExt.BusyPollRxPackets)
	netStat.TcpExtTCPAutoCorking = getValueOrDefault(curr.TcpExt.TCPAutoCorking) - getValueOrDefault(prev.TcpExt.TCPAutoCorking)
	netStat.TcpExtTCPFromZeroWindowAdv = getValueOrDefault(curr.TcpExt.TCPFromZeroWindowAdv) - getValueOrDefault(prev.TcpExt.TCPFromZeroWindowAdv)
	netStat.TcpExtTCPToZeroWindowAdv = getValueOrDefault(curr.TcpExt.TCPToZeroWindowAdv) - getValueOrDefault(prev.TcpExt.TCPToZeroWindowAdv)
	netStat.TcpExtTCPWantZeroWindowAdv = getValueOrDefault(curr.TcpExt.TCPWantZeroWindowAdv) - getValueOrDefault(prev.TcpExt.TCPWantZeroWindowAdv)
	netStat.TcpExtTCPSynRetrans = getValueOrDefault(curr.TcpExt.TCPSynRetrans) - getValueOrDefault(prev.TcpExt.TCPSynRetrans)
	netStat.TcpExtTCPOrigDataSent = getValueOrDefault(curr.TcpExt.TCPOrigDataSent) - getValueOrDefault(prev.TcpExt.TCPOrigDataSent)
	netStat.TcpExtTCPHystartTrainDetect = getValueOrDefault(curr.TcpExt.TCPHystartTrainDetect) - getValueOrDefault(prev.TcpExt.TCPHystartTrainDetect)
	netStat.TcpExtTCPHystartTrainCwnd = getValueOrDefault(curr.TcpExt.TCPHystartTrainCwnd) - getValueOrDefault(prev.TcpExt.TCPHystartTrainCwnd)
	netStat.TcpExtTCPHystartDelayDetect = getValueOrDefault(curr.TcpExt.TCPHystartDelayDetect) - getValueOrDefault(prev.TcpExt.TCPHystartDelayDetect)
	netStat.TcpExtTCPHystartDelayCwnd = getValueOrDefault(curr.TcpExt.TCPHystartDelayCwnd) - getValueOrDefault(prev.TcpExt.TCPHystartDelayCwnd)
	netStat.TcpExtTCPACKSkippedSynRecv = getValueOrDefault(curr.TcpExt.TCPACKSkippedSynRecv) - getValueOrDefault(prev.TcpExt.TCPACKSkippedSynRecv)
	netStat.TcpExtTCPACKSkippedPAWS = getValueOrDefault(curr.TcpExt.TCPACKSkippedPAWS) - getValueOrDefault(prev.TcpExt.TCPACKSkippedPAWS)
	netStat.TcpExtTCPACKSkippedSeq = getValueOrDefault(curr.TcpExt.TCPACKSkippedSeq) - getValueOrDefault(prev.TcpExt.TCPACKSkippedSeq)
	netStat.TcpExtTCPACKSkippedFinWait2 = getValueOrDefault(curr.TcpExt.TCPACKSkippedFinWait2) - getValueOrDefault(prev.TcpExt.TCPACKSkippedFinWait2)
	netStat.TcpExtTCPACKSkippedTimeWait = getValueOrDefault(curr.TcpExt.TCPACKSkippedTimeWait) - getValueOrDefault(prev.TcpExt.TCPACKSkippedTimeWait)
	netStat.TcpExtTCPACKSkippedChallenge = getValueOrDefault(curr.TcpExt.TCPACKSkippedChallenge) - getValueOrDefault(prev.TcpExt.TCPACKSkippedChallenge)
	netStat.TcpExtTCPWinProbe = getValueOrDefault(curr.TcpExt.TCPWinProbe) - getValueOrDefault(prev.TcpExt.TCPWinProbe)
	netStat.TcpExtTCPKeepAlive = getValueOrDefault(curr.TcpExt.TCPKeepAlive) - getValueOrDefault(prev.TcpExt.TCPKeepAlive)
	netStat.TcpExtTCPMTUPFail = getValueOrDefault(curr.TcpExt.TCPMTUPFail) - getValueOrDefault(prev.TcpExt.TCPMTUPFail)
	netStat.TcpExtTCPMTUPSuccess = getValueOrDefault(curr.TcpExt.TCPMTUPSuccess) - getValueOrDefault(prev.TcpExt.TCPMTUPSuccess)
	netStat.TcpExtTCPWqueueTooBig = getValueOrDefault(curr.TcpExt.TCPWqueueTooBig) - getValueOrDefault(prev.TcpExt.TCPWqueueTooBig)

	netStat.IpExtInNoRoutes = getValueOrDefault(curr.IpExt.InNoRoutes) - getValueOrDefault(prev.IpExt.InNoRoutes)
	netStat.IpExtInTruncatedPkts = getValueOrDefault(curr.IpExt.InTruncatedPkts) - getValueOrDefault(prev.IpExt.InTruncatedPkts)
	netStat.IpExtInMcastPkts = getValueOrDefault(curr.IpExt.InMcastPkts) - getValueOrDefault(prev.IpExt.InMcastPkts)
	netStat.IpExtOutMcastPkts = getValueOrDefault(curr.IpExt.OutMcastPkts) - getValueOrDefault(prev.IpExt.OutMcastPkts)
	netStat.IpExtInBcastPkts = getValueOrDefault(curr.IpExt.InBcastPkts) - getValueOrDefault(prev.IpExt.InBcastPkts)
	netStat.IpExtOutBcastPkts = getValueOrDefault(curr.IpExt.OutBcastPkts) - getValueOrDefault(prev.IpExt.OutBcastPkts)
	netStat.IpExtInOctets = getValueOrDefault(curr.IpExt.InOctets) - getValueOrDefault(prev.IpExt.InOctets)
	netStat.IpExtOutOctets = getValueOrDefault(curr.IpExt.OutOctets) - getValueOrDefault(prev.IpExt.OutOctets)
	netStat.IpExtInMcastOctets = getValueOrDefault(curr.IpExt.InMcastOctets) - getValueOrDefault(prev.IpExt.InMcastOctets)
	netStat.IpExtOutMcastOctets = getValueOrDefault(curr.IpExt.OutMcastOctets) - getValueOrDefault(prev.IpExt.OutMcastOctets)
	netStat.IpExtInBcastOctets = getValueOrDefault(curr.IpExt.InBcastOctets) - getValueOrDefault(prev.IpExt.InBcastOctets)
	netStat.IpExtOutBcastOctets = getValueOrDefault(curr.IpExt.OutBcastOctets) - getValueOrDefault(prev.IpExt.OutBcastOctets)
	netStat.IpExtInCsumErrors = getValueOrDefault(curr.IpExt.InCsumErrors) - getValueOrDefault(prev.IpExt.InCsumErrors)
	netStat.IpExtInNoECTPkts = getValueOrDefault(curr.IpExt.InNoECTPkts) - getValueOrDefault(prev.IpExt.InNoECTPkts)
	netStat.IpExtInECT1Pkts = getValueOrDefault(curr.IpExt.InECT1Pkts) - getValueOrDefault(prev.IpExt.InECT1Pkts)
	netStat.IpExtInECT0Pkts = getValueOrDefault(curr.IpExt.InECT0Pkts) - getValueOrDefault(prev.IpExt.InECT0Pkts)
	netStat.IpExtInCEPkts = getValueOrDefault(curr.IpExt.InCEPkts) - getValueOrDefault(prev.IpExt.InCEPkts)
	netStat.IpExtReasmOverlaps = getValueOrDefault(curr.IpExt.ReasmOverlaps) - getValueOrDefault(prev.IpExt.ReasmOverlaps)

}
