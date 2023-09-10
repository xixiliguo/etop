package model

import (
	"fmt"

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

func (n *NetStat) GetRenderValue(config RenderConfig, field string) string {

	s := fmt.Sprintf("no %s for network stat", field)
	switch field {
	case "IpInReceives":
		s = config[field].Render(n.IpInReceives)
	case "IpInHdrErrors":
		s = config[field].Render(n.IpInHdrErrors)
	case "IpInAddrErrors":
		s = config[field].Render(n.IpInAddrErrors)
	case "IpForwDatagrams":
		s = config[field].Render(n.IpForwDatagrams)
	case "IpInUnknownProtos":
		s = config[field].Render(n.IpInUnknownProtos)
	case "IpInDiscards":
		s = config[field].Render(n.IpInDiscards)
	case "IpInDelivers":
		s = config[field].Render(n.IpInDelivers)
	case "IpOutRequests":
		s = config[field].Render(n.IpOutRequests)
	case "IpOutDiscards":
		s = config[field].Render(n.IpOutDiscards)
	case "IpOutNoRoutes":
		s = config[field].Render(n.IpOutNoRoutes)
	case "IpReasmTimeout":
		s = config[field].Render(n.IpReasmTimeout)
	case "IpReasmReqds":
		s = config[field].Render(n.IpReasmReqds)
	case "IpReasmOKs":
		s = config[field].Render(n.IpReasmOKs)
	case "IpReasmFails":
		s = config[field].Render(n.IpReasmFails)
	case "IpFragOKs":
		s = config[field].Render(n.IpFragOKs)
	case "IpFragFails":
		s = config[field].Render(n.IpFragFails)
	case "IpFragCreates":
		s = config[field].Render(n.IpFragCreates)
	case "IcmpInMsgs":
		s = config[field].Render(n.IcmpInMsgs)
	case "IcmpInErrors":
		s = config[field].Render(n.IcmpInErrors)
	case "IcmpInCsumErrors":
		s = config[field].Render(n.IcmpInCsumErrors)
	case "IcmpInDestUnreachs":
		s = config[field].Render(n.IcmpInDestUnreachs)
	case "IcmpInTimeExcds":
		s = config[field].Render(n.IcmpInTimeExcds)
	case "IcmpInParmProbs":
		s = config[field].Render(n.IcmpInParmProbs)
	case "IcmpInSrcQuenchs":
		s = config[field].Render(n.IcmpInSrcQuenchs)
	case "IcmpInRedirects":
		s = config[field].Render(n.IcmpInRedirects)
	case "IcmpInEchos":
		s = config[field].Render(n.IcmpInEchos)
	case "IcmpInEchoReps":
		s = config[field].Render(n.IcmpInEchoReps)
	case "IcmpInTimestamps":
		s = config[field].Render(n.IcmpInTimestamps)
	case "IcmpInTimestampReps":
		s = config[field].Render(n.IcmpInTimestampReps)
	case "IcmpInAddrMasks":
		s = config[field].Render(n.IcmpInAddrMasks)
	case "IcmpInAddrMaskReps":
		s = config[field].Render(n.IcmpInAddrMaskReps)
	case "IcmpOutMsgs":
		s = config[field].Render(n.IcmpOutMsgs)
	case "IcmpOutErrors":
		s = config[field].Render(n.IcmpOutErrors)
	case "IcmpOutDestUnreachs":
		s = config[field].Render(n.IcmpOutDestUnreachs)
	case "IcmpOutTimeExcds":
		s = config[field].Render(n.IcmpOutTimeExcds)
	case "IcmpOutParmProbs":
		s = config[field].Render(n.IcmpOutParmProbs)
	case "IcmpOutSrcQuenchs":
		s = config[field].Render(n.IcmpOutSrcQuenchs)
	case "IcmpOutRedirects":
		s = config[field].Render(n.IcmpOutRedirects)
	case "IcmpOutEchos":
		s = config[field].Render(n.IcmpOutEchos)
	case "IcmpOutEchoReps":
		s = config[field].Render(n.IcmpOutEchoReps)
	case "IcmpOutTimestamps":
		s = config[field].Render(n.IcmpOutTimestamps)
	case "IcmpOutTimestampReps":
		s = config[field].Render(n.IcmpOutTimestampReps)
	case "IcmpOutAddrMasks":
		s = config[field].Render(n.IcmpOutAddrMasks)
	case "IcmpOutAddrMaskReps":
		s = config[field].Render(n.IcmpOutAddrMaskReps)
	case "IcmpInType3":
		s = config[field].Render(n.IcmpInType3)
	case "IcmpOutType3":
		s = config[field].Render(n.IcmpOutType3)
	case "TcpActiveOpens":
		s = config[field].Render(n.TcpActiveOpens)
	case "TcpPassiveOpens":
		s = config[field].Render(n.TcpPassiveOpens)
	case "TcpAttemptFails":
		s = config[field].Render(n.TcpAttemptFails)
	case "TcpEstabResets":
		s = config[field].Render(n.TcpEstabResets)
	case "TcpCurrEstab":
		s = config[field].Render(n.TcpCurrEstab)
	case "TcpInSegs":
		s = config[field].Render(n.TcpInSegs)
	case "TcpOutSegs":
		s = config[field].Render(n.TcpOutSegs)
	case "TcpRetransSegs":
		s = config[field].Render(n.TcpRetransSegs)
	case "TcpInErrs":
		s = config[field].Render(n.TcpInErrs)
	case "TcpOutRsts":
		s = config[field].Render(n.TcpOutRsts)
	case "TcpInCsumErrors":
		s = config[field].Render(n.TcpInCsumErrors)
	case "UdpInDatagrams":
		s = config[field].Render(n.UdpInDatagrams)
	case "UdpNoPorts":
		s = config[field].Render(n.UdpNoPorts)
	case "UdpInErrors":
		s = config[field].Render(n.UdpInErrors)
	case "UdpOutDatagrams":
		s = config[field].Render(n.UdpOutDatagrams)
	case "UdpRcvbufErrors":
		s = config[field].Render(n.UdpRcvbufErrors)
	case "UdpSndbufErrors":
		s = config[field].Render(n.UdpSndbufErrors)
	case "UdpInCsumErrors":
		s = config[field].Render(n.UdpInCsumErrors)
	case "UdpIgnoredMulti":
		s = config[field].Render(n.UdpIgnoredMulti)
	case "UdpLiteInDatagrams":
		s = config[field].Render(n.UdpLiteInDatagrams)
	case "UdpLiteNoPorts":
		s = config[field].Render(n.UdpLiteNoPorts)
	case "UdpLiteInErrors":
		s = config[field].Render(n.UdpLiteInErrors)
	case "UdpLiteOutDatagrams":
		s = config[field].Render(n.UdpLiteOutDatagrams)
	case "UdpLiteRcvbufErrors":
		s = config[field].Render(n.UdpLiteRcvbufErrors)
	case "UdpLiteSndbufErrors":
		s = config[field].Render(n.UdpLiteSndbufErrors)
	case "UdpLiteInCsumErrors":
		s = config[field].Render(n.UdpLiteInCsumErrors)
	case "UdpLiteIgnoredMulti":
		s = config[field].Render(n.UdpLiteIgnoredMulti)
	case "Ip6InReceives":
		s = config[field].Render(n.Ip6InReceives)
	case "Ip6InHdrErrors":
		s = config[field].Render(n.Ip6InHdrErrors)
	case "Ip6InTooBigErrors":
		s = config[field].Render(n.Ip6InTooBigErrors)
	case "Ip6InNoRoutes":
		s = config[field].Render(n.Ip6InNoRoutes)
	case "Ip6InAddrErrors":
		s = config[field].Render(n.Ip6InAddrErrors)
	case "Ip6InUnknownProtos":
		s = config[field].Render(n.Ip6InUnknownProtos)
	case "Ip6InTruncatedPkts":
		s = config[field].Render(n.Ip6InTruncatedPkts)
	case "Ip6InDiscards":
		s = config[field].Render(n.Ip6InDiscards)
	case "Ip6InDelivers":
		s = config[field].Render(n.Ip6InDelivers)
	case "Ip6OutForwDatagrams":
		s = config[field].Render(n.Ip6OutForwDatagrams)
	case "Ip6OutRequests":
		s = config[field].Render(n.Ip6OutRequests)
	case "Ip6OutDiscards":
		s = config[field].Render(n.Ip6OutDiscards)
	case "Ip6OutNoRoutes":
		s = config[field].Render(n.Ip6OutNoRoutes)
	case "Ip6ReasmTimeout":
		s = config[field].Render(n.Ip6ReasmTimeout)
	case "Ip6ReasmReqds":
		s = config[field].Render(n.Ip6ReasmReqds)
	case "Ip6ReasmOKs":
		s = config[field].Render(n.Ip6ReasmOKs)
	case "Ip6ReasmFails":
		s = config[field].Render(n.Ip6ReasmFails)
	case "Ip6FragOKs":
		s = config[field].Render(n.Ip6FragOKs)
	case "Ip6FragFails":
		s = config[field].Render(n.Ip6FragFails)
	case "Ip6FragCreates":
		s = config[field].Render(n.Ip6FragCreates)
	case "Ip6InMcastPkts":
		s = config[field].Render(n.Ip6InMcastPkts)
	case "Ip6OutMcastPkts":
		s = config[field].Render(n.Ip6OutMcastPkts)
	case "Ip6InOctets":
		s = config[field].Render(n.Ip6InOctets)
	case "Ip6OutOctets":
		s = config[field].Render(n.Ip6OutOctets)
	case "Ip6InMcastOctets":
		s = config[field].Render(n.Ip6InMcastOctets)
	case "Ip6OutMcastOctets":
		s = config[field].Render(n.Ip6OutMcastOctets)
	case "Ip6InBcastOctets":
		s = config[field].Render(n.Ip6InBcastOctets)
	case "Ip6OutBcastOctets":
		s = config[field].Render(n.Ip6OutBcastOctets)
	case "Ip6InNoECTPkts":
		s = config[field].Render(n.Ip6InNoECTPkts)
	case "Ip6InECT1Pkts":
		s = config[field].Render(n.Ip6InECT1Pkts)
	case "Ip6InECT0Pkts":
		s = config[field].Render(n.Ip6InECT0Pkts)
	case "Ip6InCEPkts":
		s = config[field].Render(n.Ip6InCEPkts)
	case "Icmp6InMsgs":
		s = config[field].Render(n.Icmp6InMsgs)
	case "Icmp6InErrors":
		s = config[field].Render(n.Icmp6InErrors)
	case "Icmp6OutMsgs":
		s = config[field].Render(n.Icmp6OutMsgs)
	case "Icmp6OutErrors":
		s = config[field].Render(n.Icmp6OutErrors)
	case "Icmp6InCsumErrors":
		s = config[field].Render(n.Icmp6InCsumErrors)
	case "Icmp6InDestUnreachs":
		s = config[field].Render(n.Icmp6InDestUnreachs)
	case "Icmp6InPktTooBigs":
		s = config[field].Render(n.Icmp6InPktTooBigs)
	case "Icmp6InTimeExcds":
		s = config[field].Render(n.Icmp6InTimeExcds)
	case "Icmp6InParmProblems":
		s = config[field].Render(n.Icmp6InParmProblems)
	case "Icmp6InEchos":
		s = config[field].Render(n.Icmp6InEchos)
	case "Icmp6InEchoReplies":
		s = config[field].Render(n.Icmp6InEchoReplies)
	case "Icmp6InGroupMembQueries":
		s = config[field].Render(n.Icmp6InGroupMembQueries)
	case "Icmp6InGroupMembResponses":
		s = config[field].Render(n.Icmp6InGroupMembResponses)
	case "Icmp6InGroupMembReductions":
		s = config[field].Render(n.Icmp6InGroupMembReductions)
	case "Icmp6InRouterSolicits":
		s = config[field].Render(n.Icmp6InRouterSolicits)
	case "Icmp6InRouterAdvertisements":
		s = config[field].Render(n.Icmp6InRouterAdvertisements)
	case "Icmp6InNeighborSolicits":
		s = config[field].Render(n.Icmp6InNeighborSolicits)
	case "Icmp6InNeighborAdvertisements":
		s = config[field].Render(n.Icmp6InNeighborAdvertisements)
	case "Icmp6InRedirects":
		s = config[field].Render(n.Icmp6InRedirects)
	case "Icmp6InMLDv2Reports":
		s = config[field].Render(n.Icmp6InMLDv2Reports)
	case "Icmp6OutDestUnreachs":
		s = config[field].Render(n.Icmp6OutDestUnreachs)
	case "Icmp6OutPktTooBigs":
		s = config[field].Render(n.Icmp6OutPktTooBigs)
	case "Icmp6OutTimeExcds":
		s = config[field].Render(n.Icmp6OutTimeExcds)
	case "Icmp6OutParmProblems":
		s = config[field].Render(n.Icmp6OutParmProblems)
	case "Icmp6OutEchos":
		s = config[field].Render(n.Icmp6OutEchos)
	case "Icmp6OutEchoReplies":
		s = config[field].Render(n.Icmp6OutEchoReplies)
	case "Icmp6OutGroupMembQueries":
		s = config[field].Render(n.Icmp6OutGroupMembQueries)
	case "Icmp6OutGroupMembResponses":
		s = config[field].Render(n.Icmp6OutGroupMembResponses)
	case "Icmp6OutGroupMembReductions":
		s = config[field].Render(n.Icmp6OutGroupMembReductions)
	case "Icmp6OutRouterSolicits":
		s = config[field].Render(n.Icmp6OutRouterSolicits)
	case "Icmp6OutRouterAdvertisements":
		s = config[field].Render(n.Icmp6OutRouterAdvertisements)
	case "Icmp6OutNeighborSolicits":
		s = config[field].Render(n.Icmp6OutNeighborSolicits)
	case "Icmp6OutNeighborAdvertisements":
		s = config[field].Render(n.Icmp6OutNeighborAdvertisements)
	case "Icmp6OutRedirects":
		s = config[field].Render(n.Icmp6OutRedirects)
	case "Icmp6OutMLDv2Reports":
		s = config[field].Render(n.Icmp6OutMLDv2Reports)
	case "Icmp6InType1":
		s = config[field].Render(n.Icmp6InType1)
	case "Icmp6InType134":
		s = config[field].Render(n.Icmp6InType134)
	case "Icmp6InType135":
		s = config[field].Render(n.Icmp6InType135)
	case "Icmp6InType136":
		s = config[field].Render(n.Icmp6InType136)
	case "Icmp6InType143":
		s = config[field].Render(n.Icmp6InType143)
	case "Icmp6OutType133":
		s = config[field].Render(n.Icmp6OutType133)
	case "Icmp6OutType135":
		s = config[field].Render(n.Icmp6OutType135)
	case "Icmp6OutType136":
		s = config[field].Render(n.Icmp6OutType136)
	case "Icmp6OutType143":
		s = config[field].Render(n.Icmp6OutType143)
	case "Udp6InDatagrams":
		s = config[field].Render(n.Udp6InDatagrams)
	case "Udp6NoPorts":
		s = config[field].Render(n.Udp6NoPorts)
	case "Udp6InErrors":
		s = config[field].Render(n.Udp6InErrors)
	case "Udp6OutDatagrams":
		s = config[field].Render(n.Udp6OutDatagrams)
	case "Udp6RcvbufErrors":
		s = config[field].Render(n.Udp6RcvbufErrors)
	case "Udp6SndbufErrors":
		s = config[field].Render(n.Udp6SndbufErrors)
	case "Udp6InCsumErrors":
		s = config[field].Render(n.Udp6InCsumErrors)
	case "Udp6IgnoredMulti":
		s = config[field].Render(n.Udp6IgnoredMulti)
	case "UdpLite6InDatagrams":
		s = config[field].Render(n.UdpLite6InDatagrams)
	case "UdpLite6NoPorts":
		s = config[field].Render(n.UdpLite6NoPorts)
	case "UdpLite6InErrors":
		s = config[field].Render(n.UdpLite6InErrors)
	case "UdpLite6OutDatagrams":
		s = config[field].Render(n.UdpLite6OutDatagrams)
	case "UdpLite6RcvbufErrors":
		s = config[field].Render(n.UdpLite6RcvbufErrors)
	case "UdpLite6SndbufErrors":
		s = config[field].Render(n.UdpLite6SndbufErrors)
	case "UdpLite6InCsumErrors":
		s = config[field].Render(n.UdpLite6InCsumErrors)
	case "TcpExtSyncookiesSent":
		s = config[field].Render(n.TcpExtSyncookiesSent)
	case "TcpExtSyncookiesRecv":
		s = config[field].Render(n.TcpExtSyncookiesRecv)
	case "TcpExtSyncookiesFailed":
		s = config[field].Render(n.TcpExtSyncookiesFailed)
	case "TcpExtEmbryonicRsts":
		s = config[field].Render(n.TcpExtEmbryonicRsts)
	case "TcpExtPruneCalled":
		s = config[field].Render(n.TcpExtPruneCalled)
	case "TcpExtRcvPruned":
		s = config[field].Render(n.TcpExtRcvPruned)
	case "TcpExtOfoPruned":
		s = config[field].Render(n.TcpExtOfoPruned)
	case "TcpExtOutOfWindowIcmps":
		s = config[field].Render(n.TcpExtOutOfWindowIcmps)
	case "TcpExtLockDroppedIcmps":
		s = config[field].Render(n.TcpExtLockDroppedIcmps)
	case "TcpExtArpFilter":
		s = config[field].Render(n.TcpExtArpFilter)
	case "TcpExtTW":
		s = config[field].Render(n.TcpExtTW)
	case "TcpExtTWRecycled":
		s = config[field].Render(n.TcpExtTWRecycled)
	case "TcpExtTWKilled":
		s = config[field].Render(n.TcpExtTWKilled)
	case "TcpExtPAWSActive":
		s = config[field].Render(n.TcpExtPAWSActive)
	case "TcpExtPAWSEstab":
		s = config[field].Render(n.TcpExtPAWSEstab)
	case "TcpExtDelayedACKs":
		s = config[field].Render(n.TcpExtDelayedACKs)
	case "TcpExtDelayedACKLocked":
		s = config[field].Render(n.TcpExtDelayedACKLocked)
	case "TcpExtDelayedACKLost":
		s = config[field].Render(n.TcpExtDelayedACKLost)
	case "TcpExtListenOverflows":
		s = config[field].Render(n.TcpExtListenOverflows)
	case "TcpExtListenDrops":
		s = config[field].Render(n.TcpExtListenDrops)
	case "TcpExtTCPHPHits":
		s = config[field].Render(n.TcpExtTCPHPHits)
	case "TcpExtTCPPureAcks":
		s = config[field].Render(n.TcpExtTCPPureAcks)
	case "TcpExtTCPHPAcks":
		s = config[field].Render(n.TcpExtTCPHPAcks)
	case "TcpExtTCPRenoRecovery":
		s = config[field].Render(n.TcpExtTCPRenoRecovery)
	case "TcpExtTCPSackRecovery":
		s = config[field].Render(n.TcpExtTCPSackRecovery)
	case "TcpExtTCPSACKReneging":
		s = config[field].Render(n.TcpExtTCPSACKReneging)
	case "TcpExtTCPSACKReorder":
		s = config[field].Render(n.TcpExtTCPSACKReorder)
	case "TcpExtTCPRenoReorder":
		s = config[field].Render(n.TcpExtTCPRenoReorder)
	case "TcpExtTCPTSReorder":
		s = config[field].Render(n.TcpExtTCPTSReorder)
	case "TcpExtTCPFullUndo":
		s = config[field].Render(n.TcpExtTCPFullUndo)
	case "TcpExtTCPPartialUndo":
		s = config[field].Render(n.TcpExtTCPPartialUndo)
	case "TcpExtTCPDSACKUndo":
		s = config[field].Render(n.TcpExtTCPDSACKUndo)
	case "TcpExtTCPLossUndo":
		s = config[field].Render(n.TcpExtTCPLossUndo)
	case "TcpExtTCPLostRetransmit":
		s = config[field].Render(n.TcpExtTCPLostRetransmit)
	case "TcpExtTCPRenoFailures":
		s = config[field].Render(n.TcpExtTCPRenoFailures)
	case "TcpExtTCPSackFailures":
		s = config[field].Render(n.TcpExtTCPSackFailures)
	case "TcpExtTCPLossFailures":
		s = config[field].Render(n.TcpExtTCPLossFailures)
	case "TcpExtTCPFastRetrans":
		s = config[field].Render(n.TcpExtTCPFastRetrans)
	case "TcpExtTCPSlowStartRetrans":
		s = config[field].Render(n.TcpExtTCPSlowStartRetrans)
	case "TcpExtTCPTimeouts":
		s = config[field].Render(n.TcpExtTCPTimeouts)
	case "TcpExtTCPLossProbes":
		s = config[field].Render(n.TcpExtTCPLossProbes)
	case "TcpExtTCPLossProbeRecovery":
		s = config[field].Render(n.TcpExtTCPLossProbeRecovery)
	case "TcpExtTCPRenoRecoveryFail":
		s = config[field].Render(n.TcpExtTCPRenoRecoveryFail)
	case "TcpExtTCPSackRecoveryFail":
		s = config[field].Render(n.TcpExtTCPSackRecoveryFail)
	case "TcpExtTCPRcvCollapsed":
		s = config[field].Render(n.TcpExtTCPRcvCollapsed)
	case "TcpExtTCPDSACKOldSent":
		s = config[field].Render(n.TcpExtTCPDSACKOldSent)
	case "TcpExtTCPDSACKOfoSent":
		s = config[field].Render(n.TcpExtTCPDSACKOfoSent)
	case "TcpExtTCPDSACKRecv":
		s = config[field].Render(n.TcpExtTCPDSACKRecv)
	case "TcpExtTCPDSACKOfoRecv":
		s = config[field].Render(n.TcpExtTCPDSACKOfoRecv)
	case "TcpExtTCPAbortOnData":
		s = config[field].Render(n.TcpExtTCPAbortOnData)
	case "TcpExtTCPAbortOnClose":
		s = config[field].Render(n.TcpExtTCPAbortOnClose)
	case "TcpExtTCPAbortOnMemory":
		s = config[field].Render(n.TcpExtTCPAbortOnMemory)
	case "TcpExtTCPAbortOnTimeout":
		s = config[field].Render(n.TcpExtTCPAbortOnTimeout)
	case "TcpExtTCPAbortOnLinger":
		s = config[field].Render(n.TcpExtTCPAbortOnLinger)
	case "TcpExtTCPAbortFailed":
		s = config[field].Render(n.TcpExtTCPAbortFailed)
	case "TcpExtTCPMemoryPressures":
		s = config[field].Render(n.TcpExtTCPMemoryPressures)
	case "TcpExtTCPMemoryPressuresChrono":
		s = config[field].Render(n.TcpExtTCPMemoryPressuresChrono)
	case "TcpExtTCPSACKDiscard":
		s = config[field].Render(n.TcpExtTCPSACKDiscard)
	case "TcpExtTCPDSACKIgnoredOld":
		s = config[field].Render(n.TcpExtTCPDSACKIgnoredOld)
	case "TcpExtTCPDSACKIgnoredNoUndo":
		s = config[field].Render(n.TcpExtTCPDSACKIgnoredNoUndo)
	case "TcpExtTCPSpuriousRTOs":
		s = config[field].Render(n.TcpExtTCPSpuriousRTOs)
	case "TcpExtTCPMD5NotFound":
		s = config[field].Render(n.TcpExtTCPMD5NotFound)
	case "TcpExtTCPMD5Unexpected":
		s = config[field].Render(n.TcpExtTCPMD5Unexpected)
	case "TcpExtTCPMD5Failure":
		s = config[field].Render(n.TcpExtTCPMD5Failure)
	case "TcpExtTCPSackShifted":
		s = config[field].Render(n.TcpExtTCPSackShifted)
	case "TcpExtTCPSackMerged":
		s = config[field].Render(n.TcpExtTCPSackMerged)
	case "TcpExtTCPSackShiftFallback":
		s = config[field].Render(n.TcpExtTCPSackShiftFallback)
	case "TcpExtTCPBacklogDrop":
		s = config[field].Render(n.TcpExtTCPBacklogDrop)
	case "TcpExtPFMemallocDrop":
		s = config[field].Render(n.TcpExtPFMemallocDrop)
	case "TcpExtTCPMinTTLDrop":
		s = config[field].Render(n.TcpExtTCPMinTTLDrop)
	case "TcpExtTCPDeferAcceptDrop":
		s = config[field].Render(n.TcpExtTCPDeferAcceptDrop)
	case "TcpExtIPReversePathFilter":
		s = config[field].Render(n.TcpExtIPReversePathFilter)
	case "TcpExtTCPTimeWaitOverflow":
		s = config[field].Render(n.TcpExtTCPTimeWaitOverflow)
	case "TcpExtTCPReqQFullDoCookies":
		s = config[field].Render(n.TcpExtTCPReqQFullDoCookies)
	case "TcpExtTCPReqQFullDrop":
		s = config[field].Render(n.TcpExtTCPReqQFullDrop)
	case "TcpExtTCPRetransFail":
		s = config[field].Render(n.TcpExtTCPRetransFail)
	case "TcpExtTCPRcvCoalesce":
		s = config[field].Render(n.TcpExtTCPRcvCoalesce)
	case "TcpExtTCPRcvQDrop":
		s = config[field].Render(n.TcpExtTCPRcvQDrop)
	case "TcpExtTCPOFOQueue":
		s = config[field].Render(n.TcpExtTCPOFOQueue)
	case "TcpExtTCPOFODrop":
		s = config[field].Render(n.TcpExtTCPOFODrop)
	case "TcpExtTCPOFOMerge":
		s = config[field].Render(n.TcpExtTCPOFOMerge)
	case "TcpExtTCPChallengeACK":
		s = config[field].Render(n.TcpExtTCPChallengeACK)
	case "TcpExtTCPSYNChallenge":
		s = config[field].Render(n.TcpExtTCPSYNChallenge)
	case "TcpExtTCPFastOpenActive":
		s = config[field].Render(n.TcpExtTCPFastOpenActive)
	case "TcpExtTCPFastOpenActiveFail":
		s = config[field].Render(n.TcpExtTCPFastOpenActiveFail)
	case "TcpExtTCPFastOpenPassive":
		s = config[field].Render(n.TcpExtTCPFastOpenPassive)
	case "TcpExtTCPFastOpenPassiveFail":
		s = config[field].Render(n.TcpExtTCPFastOpenPassiveFail)
	case "TcpExtTCPFastOpenListenOverflow":
		s = config[field].Render(n.TcpExtTCPFastOpenListenOverflow)
	case "TcpExtTCPFastOpenCookieReqd":
		s = config[field].Render(n.TcpExtTCPFastOpenCookieReqd)
	case "TcpExtTCPFastOpenBlackhole":
		s = config[field].Render(n.TcpExtTCPFastOpenBlackhole)
	case "TcpExtTCPSpuriousRtxHostQueues":
		s = config[field].Render(n.TcpExtTCPSpuriousRtxHostQueues)
	case "TcpExtBusyPollRxPackets":
		s = config[field].Render(n.TcpExtBusyPollRxPackets)
	case "TcpExtTCPAutoCorking":
		s = config[field].Render(n.TcpExtTCPAutoCorking)
	case "TcpExtTCPFromZeroWindowAdv":
		s = config[field].Render(n.TcpExtTCPFromZeroWindowAdv)
	case "TcpExtTCPToZeroWindowAdv":
		s = config[field].Render(n.TcpExtTCPToZeroWindowAdv)
	case "TcpExtTCPWantZeroWindowAdv":
		s = config[field].Render(n.TcpExtTCPWantZeroWindowAdv)
	case "TcpExtTCPSynRetrans":
		s = config[field].Render(n.TcpExtTCPSynRetrans)
	case "TcpExtTCPOrigDataSent":
		s = config[field].Render(n.TcpExtTCPOrigDataSent)
	case "TcpExtTCPHystartTrainDetect":
		s = config[field].Render(n.TcpExtTCPHystartTrainDetect)
	case "TcpExtTCPHystartTrainCwnd":
		s = config[field].Render(n.TcpExtTCPHystartTrainCwnd)
	case "TcpExtTCPHystartDelayDetect":
		s = config[field].Render(n.TcpExtTCPHystartDelayDetect)
	case "TcpExtTCPHystartDelayCwnd":
		s = config[field].Render(n.TcpExtTCPHystartDelayCwnd)
	case "TcpExtTCPACKSkippedSynRecv":
		s = config[field].Render(n.TcpExtTCPACKSkippedSynRecv)
	case "TcpExtTCPACKSkippedPAWS":
		s = config[field].Render(n.TcpExtTCPACKSkippedPAWS)
	case "TcpExtTCPACKSkippedSeq":
		s = config[field].Render(n.TcpExtTCPACKSkippedSeq)
	case "TcpExtTCPACKSkippedFinWait2":
		s = config[field].Render(n.TcpExtTCPACKSkippedFinWait2)
	case "TcpExtTCPACKSkippedTimeWait":
		s = config[field].Render(n.TcpExtTCPACKSkippedTimeWait)
	case "TcpExtTCPACKSkippedChallenge":
		s = config[field].Render(n.TcpExtTCPACKSkippedChallenge)
	case "TcpExtTCPWinProbe":
		s = config[field].Render(n.TcpExtTCPWinProbe)
	case "TcpExtTCPKeepAlive":
		s = config[field].Render(n.TcpExtTCPKeepAlive)
	case "TcpExtTCPMTUPFail":
		s = config[field].Render(n.TcpExtTCPMTUPFail)
	case "TcpExtTCPMTUPSuccess":
		s = config[field].Render(n.TcpExtTCPMTUPSuccess)
	case "TcpExtTCPWqueueTooBig":
		s = config[field].Render(n.TcpExtTCPWqueueTooBig)
	case "IpExtInNoRoutes":
		s = config[field].Render(n.IpExtInNoRoutes)
	case "IpExtInTruncatedPkts":
		s = config[field].Render(n.IpExtInTruncatedPkts)
	case "IpExtInMcastPkts":
		s = config[field].Render(n.IpExtInMcastPkts)
	case "IpExtOutMcastPkts":
		s = config[field].Render(n.IpExtOutMcastPkts)
	case "IpExtInBcastPkts":
		s = config[field].Render(n.IpExtInBcastPkts)
	case "IpExtOutBcastPkts":
		s = config[field].Render(n.IpExtOutBcastPkts)
	case "IpExtInOctets":
		s = config[field].Render(n.IpExtInOctets)
	case "IpExtOutOctets":
		s = config[field].Render(n.IpExtOutOctets)
	case "IpExtInMcastOctets":
		s = config[field].Render(n.IpExtInMcastOctets)
	case "IpExtOutMcastOctets":
		s = config[field].Render(n.IpExtOutMcastOctets)
	case "IpExtInBcastOctets":
		s = config[field].Render(n.IpExtInBcastOctets)
	case "IpExtOutBcastOctets":
		s = config[field].Render(n.IpExtOutBcastOctets)
	case "IpExtInCsumErrors":
		s = config[field].Render(n.IpExtInCsumErrors)
	case "IpExtInNoECTPkts":
		s = config[field].Render(n.IpExtInNoECTPkts)
	case "IpExtInECT1Pkts":
		s = config[field].Render(n.IpExtInECT1Pkts)
	case "IpExtInECT0Pkts":
		s = config[field].Render(n.IpExtInECT0Pkts)
	case "IpExtInCEPkts":
		s = config[field].Render(n.IpExtInCEPkts)
	case "IpExtReasmOverlaps":
		s = config[field].Render(n.IpExtReasmOverlaps)
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
