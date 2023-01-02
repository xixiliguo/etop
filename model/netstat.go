package model

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

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

	netStat.IpInReceives = getValueOrDefault(curr.ProcSnmp.InReceives) - getValueOrDefault(prev.ProcSnmp.InReceives)
	netStat.IpInHdrErrors = getValueOrDefault(curr.ProcSnmp.InHdrErrors) - getValueOrDefault(prev.ProcSnmp.InHdrErrors)
	netStat.IpInAddrErrors = getValueOrDefault(curr.ProcSnmp.InAddrErrors) - getValueOrDefault(prev.ProcSnmp.InAddrErrors)
	netStat.IpForwDatagrams = getValueOrDefault(curr.ProcSnmp.ForwDatagrams) - getValueOrDefault(prev.ProcSnmp.ForwDatagrams)
	netStat.IpInUnknownProtos = getValueOrDefault(curr.ProcSnmp.InUnknownProtos) - getValueOrDefault(prev.ProcSnmp.InUnknownProtos)
	netStat.IpInDiscards = getValueOrDefault(curr.ProcSnmp.InDiscards) - getValueOrDefault(prev.ProcSnmp.InDiscards)
	netStat.IpInDelivers = getValueOrDefault(curr.ProcSnmp.InDelivers) - getValueOrDefault(prev.ProcSnmp.InDelivers)
	netStat.IpOutRequests = getValueOrDefault(curr.ProcSnmp.OutRequests) - getValueOrDefault(prev.ProcSnmp.OutRequests)
	netStat.IpOutDiscards = getValueOrDefault(curr.ProcSnmp.OutDiscards) - getValueOrDefault(prev.ProcSnmp.OutDiscards)
	netStat.IpOutNoRoutes = getValueOrDefault(curr.ProcSnmp.OutNoRoutes) - getValueOrDefault(prev.ProcSnmp.OutNoRoutes)
	netStat.IpReasmTimeout = getValueOrDefault(curr.ProcSnmp.ReasmTimeout) - getValueOrDefault(prev.ProcSnmp.ReasmTimeout)
	netStat.IpReasmReqds = getValueOrDefault(curr.ProcSnmp.ReasmReqds) - getValueOrDefault(prev.ProcSnmp.ReasmReqds)
	netStat.IpReasmOKs = getValueOrDefault(curr.ProcSnmp.ReasmOKs) - getValueOrDefault(prev.ProcSnmp.ReasmOKs)
	netStat.IpReasmFails = getValueOrDefault(curr.ProcSnmp.ReasmFails) - getValueOrDefault(prev.ProcSnmp.ReasmFails)
	netStat.IpFragOKs = getValueOrDefault(curr.ProcSnmp.FragOKs) - getValueOrDefault(prev.ProcSnmp.FragOKs)
	netStat.IpFragFails = getValueOrDefault(curr.ProcSnmp.FragFails) - getValueOrDefault(prev.ProcSnmp.FragFails)
	netStat.IpFragCreates = getValueOrDefault(curr.ProcSnmp.FragCreates) - getValueOrDefault(prev.ProcSnmp.FragCreates)

	netStat.IcmpInMsgs = getValueOrDefault(curr.ProcSnmp.InMsgs) - getValueOrDefault(prev.ProcSnmp.InMsgs)
	netStat.IcmpInErrors = getValueOrDefault(curr.ProcSnmp.Icmp.InErrors) - getValueOrDefault(prev.ProcSnmp.Icmp.InErrors)
	netStat.IcmpInCsumErrors = getValueOrDefault(curr.ProcSnmp.Icmp.InCsumErrors) - getValueOrDefault(prev.ProcSnmp.Icmp.InCsumErrors)
	netStat.IcmpInDestUnreachs = getValueOrDefault(curr.ProcSnmp.InDestUnreachs) - getValueOrDefault(prev.ProcSnmp.InDestUnreachs)
	netStat.IcmpInTimeExcds = getValueOrDefault(curr.ProcSnmp.InTimeExcds) - getValueOrDefault(prev.ProcSnmp.InTimeExcds)
	netStat.IcmpInParmProbs = getValueOrDefault(curr.ProcSnmp.InParmProbs) - getValueOrDefault(prev.ProcSnmp.InParmProbs)
	netStat.IcmpInSrcQuenchs = getValueOrDefault(curr.ProcSnmp.InSrcQuenchs) - getValueOrDefault(prev.ProcSnmp.InSrcQuenchs)
	netStat.IcmpInRedirects = getValueOrDefault(curr.ProcSnmp.InRedirects) - getValueOrDefault(prev.ProcSnmp.InRedirects)
	netStat.IcmpInEchos = getValueOrDefault(curr.ProcSnmp.InEchos) - getValueOrDefault(prev.ProcSnmp.InEchos)
	netStat.IcmpInEchoReps = getValueOrDefault(curr.ProcSnmp.InEchoReps) - getValueOrDefault(prev.ProcSnmp.InEchoReps)
	netStat.IcmpInTimestamps = getValueOrDefault(curr.ProcSnmp.InTimestamps) - getValueOrDefault(prev.ProcSnmp.InTimestamps)
	netStat.IcmpInTimestampReps = getValueOrDefault(curr.ProcSnmp.InTimestampReps) - getValueOrDefault(prev.ProcSnmp.InTimestampReps)
	netStat.IcmpInAddrMasks = getValueOrDefault(curr.ProcSnmp.InAddrMasks) - getValueOrDefault(prev.ProcSnmp.InAddrMasks)
	netStat.IcmpInAddrMaskReps = getValueOrDefault(curr.ProcSnmp.InAddrMaskReps) - getValueOrDefault(prev.ProcSnmp.InAddrMaskReps)
	netStat.IcmpOutMsgs = getValueOrDefault(curr.ProcSnmp.OutMsgs) - getValueOrDefault(prev.ProcSnmp.OutMsgs)
	netStat.IcmpOutErrors = getValueOrDefault(curr.ProcSnmp.OutErrors) - getValueOrDefault(prev.ProcSnmp.OutErrors)
	netStat.IcmpOutDestUnreachs = getValueOrDefault(curr.ProcSnmp.OutDestUnreachs) - getValueOrDefault(prev.ProcSnmp.OutDestUnreachs)
	netStat.IcmpOutTimeExcds = getValueOrDefault(curr.ProcSnmp.OutTimeExcds) - getValueOrDefault(prev.ProcSnmp.OutTimeExcds)
	netStat.IcmpOutParmProbs = getValueOrDefault(curr.ProcSnmp.OutParmProbs) - getValueOrDefault(prev.ProcSnmp.OutParmProbs)
	netStat.IcmpOutSrcQuenchs = getValueOrDefault(curr.ProcSnmp.OutSrcQuenchs) - getValueOrDefault(prev.ProcSnmp.OutSrcQuenchs)
	netStat.IcmpOutRedirects = getValueOrDefault(curr.ProcSnmp.OutRedirects) - getValueOrDefault(prev.ProcSnmp.OutRedirects)
	netStat.IcmpOutEchos = getValueOrDefault(curr.ProcSnmp.OutEchos) - getValueOrDefault(prev.ProcSnmp.OutEchos)
	netStat.IcmpOutEchoReps = getValueOrDefault(curr.ProcSnmp.OutEchoReps) - getValueOrDefault(prev.ProcSnmp.OutEchoReps)
	netStat.IcmpOutTimestamps = getValueOrDefault(curr.ProcSnmp.OutTimestamps) - getValueOrDefault(prev.ProcSnmp.OutTimestamps)
	netStat.IcmpOutTimestampReps = getValueOrDefault(curr.ProcSnmp.OutTimestampReps) - getValueOrDefault(prev.ProcSnmp.OutTimestampReps)
	netStat.IcmpOutAddrMasks = getValueOrDefault(curr.ProcSnmp.OutAddrMasks) - getValueOrDefault(prev.ProcSnmp.OutAddrMasks)
	netStat.IcmpOutAddrMaskReps = getValueOrDefault(curr.ProcSnmp.OutAddrMaskReps) - getValueOrDefault(prev.ProcSnmp.OutAddrMaskReps)
	netStat.IcmpInType3 = getValueOrDefault(curr.ProcSnmp.InType3) - getValueOrDefault(prev.ProcSnmp.InType3)
	netStat.IcmpOutType3 = getValueOrDefault(curr.ProcSnmp.OutType3) - getValueOrDefault(prev.ProcSnmp.OutType3)

	netStat.TcpActiveOpens = getValueOrDefault(curr.ProcSnmp.ActiveOpens) - getValueOrDefault(prev.ProcSnmp.ActiveOpens)
	netStat.TcpPassiveOpens = getValueOrDefault(curr.ProcSnmp.PassiveOpens) - getValueOrDefault(prev.ProcSnmp.PassiveOpens)
	netStat.TcpAttemptFails = getValueOrDefault(curr.ProcSnmp.AttemptFails) - getValueOrDefault(prev.ProcSnmp.AttemptFails)
	netStat.TcpEstabResets = getValueOrDefault(curr.ProcSnmp.EstabResets) - getValueOrDefault(prev.ProcSnmp.EstabResets)
	netStat.TcpCurrEstab = getValueOrDefault(curr.ProcSnmp.CurrEstab)
	netStat.TcpInSegs = getValueOrDefault(curr.ProcSnmp.InSegs) - getValueOrDefault(prev.ProcSnmp.InSegs)
	netStat.TcpOutSegs = getValueOrDefault(curr.ProcSnmp.OutSegs) - getValueOrDefault(prev.ProcSnmp.OutSegs)
	netStat.TcpRetransSegs = getValueOrDefault(curr.ProcSnmp.RetransSegs) - getValueOrDefault(prev.ProcSnmp.RetransSegs)
	netStat.TcpInErrs = getValueOrDefault(curr.ProcSnmp.InErrs) - getValueOrDefault(prev.ProcSnmp.InErrs)
	netStat.TcpOutRsts = getValueOrDefault(curr.ProcSnmp.OutRsts) - getValueOrDefault(prev.ProcSnmp.OutRsts)
	netStat.TcpInCsumErrors = getValueOrDefault(curr.ProcSnmp.Tcp.InCsumErrors) - getValueOrDefault(prev.ProcSnmp.Tcp.InCsumErrors)

	netStat.UdpInDatagrams = getValueOrDefault(curr.ProcSnmp.Udp.InDatagrams) - getValueOrDefault(prev.ProcSnmp.Udp.InDatagrams)
	netStat.UdpNoPorts = getValueOrDefault(curr.ProcSnmp.Udp.NoPorts) - getValueOrDefault(prev.ProcSnmp.Udp.NoPorts)
	netStat.UdpInErrors = getValueOrDefault(curr.ProcSnmp.Udp.InErrors) - getValueOrDefault(prev.ProcSnmp.Udp.InErrors)
	netStat.UdpOutDatagrams = getValueOrDefault(curr.ProcSnmp.Udp.OutDatagrams) - getValueOrDefault(prev.ProcSnmp.Udp.OutDatagrams)
	netStat.UdpRcvbufErrors = getValueOrDefault(curr.ProcSnmp.Udp.RcvbufErrors) - getValueOrDefault(prev.ProcSnmp.Udp.RcvbufErrors)
	netStat.UdpSndbufErrors = getValueOrDefault(curr.ProcSnmp.Udp.SndbufErrors) - getValueOrDefault(prev.ProcSnmp.Udp.SndbufErrors)
	netStat.UdpInCsumErrors = getValueOrDefault(curr.ProcSnmp.Udp.InCsumErrors) - getValueOrDefault(prev.ProcSnmp.Udp.InCsumErrors)
	netStat.UdpIgnoredMulti = getValueOrDefault(curr.ProcSnmp.Udp.IgnoredMulti) - getValueOrDefault(prev.ProcSnmp.Udp.IgnoredMulti)

	netStat.UdpLiteInDatagrams = getValueOrDefault(curr.ProcSnmp.UdpLite.InDatagrams) - getValueOrDefault(prev.ProcSnmp.UdpLite.InDatagrams)
	netStat.UdpLiteNoPorts = getValueOrDefault(curr.ProcSnmp.UdpLite.NoPorts) - getValueOrDefault(prev.ProcSnmp.UdpLite.NoPorts)
	netStat.UdpLiteInErrors = getValueOrDefault(curr.ProcSnmp.UdpLite.InErrors) - getValueOrDefault(prev.ProcSnmp.UdpLite.InErrors)
	netStat.UdpLiteOutDatagrams = getValueOrDefault(curr.ProcSnmp.UdpLite.OutDatagrams) - getValueOrDefault(prev.ProcSnmp.UdpLite.OutDatagrams)
	netStat.UdpLiteRcvbufErrors = getValueOrDefault(curr.ProcSnmp.UdpLite.RcvbufErrors) - getValueOrDefault(prev.ProcSnmp.UdpLite.RcvbufErrors)
	netStat.UdpLiteSndbufErrors = getValueOrDefault(curr.ProcSnmp.UdpLite.SndbufErrors) - getValueOrDefault(prev.ProcSnmp.UdpLite.SndbufErrors)
	netStat.UdpLiteInCsumErrors = getValueOrDefault(curr.ProcSnmp.UdpLite.InCsumErrors) - getValueOrDefault(prev.ProcSnmp.UdpLite.InCsumErrors)
	netStat.UdpLiteIgnoredMulti = getValueOrDefault(curr.ProcSnmp.UdpLite.IgnoredMulti) - getValueOrDefault(prev.ProcSnmp.UdpLite.IgnoredMulti)

	netStat.Ip6InReceives = getValueOrDefault(curr.ProcSnmp6.InReceives) - getValueOrDefault(prev.ProcSnmp6.InReceives)
	netStat.Ip6InHdrErrors = getValueOrDefault(curr.ProcSnmp6.InHdrErrors) - getValueOrDefault(prev.ProcSnmp6.InHdrErrors)
	netStat.Ip6InTooBigErrors = getValueOrDefault(curr.ProcSnmp6.InTooBigErrors) - getValueOrDefault(prev.ProcSnmp6.InTooBigErrors)
	netStat.Ip6InNoRoutes = getValueOrDefault(curr.ProcSnmp6.InNoRoutes) - getValueOrDefault(prev.ProcSnmp6.InNoRoutes)
	netStat.Ip6InAddrErrors = getValueOrDefault(curr.ProcSnmp6.InAddrErrors) - getValueOrDefault(prev.ProcSnmp6.InAddrErrors)
	netStat.Ip6InUnknownProtos = getValueOrDefault(curr.ProcSnmp6.InUnknownProtos) - getValueOrDefault(prev.ProcSnmp6.InUnknownProtos)
	netStat.Ip6InTruncatedPkts = getValueOrDefault(curr.ProcSnmp6.InTruncatedPkts) - getValueOrDefault(prev.ProcSnmp6.InTruncatedPkts)
	netStat.Ip6InDiscards = getValueOrDefault(curr.ProcSnmp6.InDiscards) - getValueOrDefault(prev.ProcSnmp6.InDiscards)
	netStat.Ip6InDelivers = getValueOrDefault(curr.ProcSnmp6.InDelivers) - getValueOrDefault(prev.ProcSnmp6.InDelivers)
	netStat.Ip6OutForwDatagrams = getValueOrDefault(curr.ProcSnmp6.OutForwDatagrams) - getValueOrDefault(prev.ProcSnmp6.OutForwDatagrams)
	netStat.Ip6OutRequests = getValueOrDefault(curr.ProcSnmp6.OutRequests) - getValueOrDefault(prev.ProcSnmp6.OutRequests)
	netStat.Ip6OutDiscards = getValueOrDefault(curr.ProcSnmp6.OutDiscards) - getValueOrDefault(prev.ProcSnmp6.OutDiscards)
	netStat.Ip6OutNoRoutes = getValueOrDefault(curr.ProcSnmp6.OutNoRoutes) - getValueOrDefault(prev.ProcSnmp6.OutNoRoutes)
	netStat.Ip6ReasmTimeout = getValueOrDefault(curr.ProcSnmp6.ReasmTimeout) - getValueOrDefault(prev.ProcSnmp6.ReasmTimeout)
	netStat.Ip6ReasmReqds = getValueOrDefault(curr.ProcSnmp6.ReasmReqds) - getValueOrDefault(prev.ProcSnmp6.ReasmReqds)
	netStat.Ip6ReasmOKs = getValueOrDefault(curr.ProcSnmp6.ReasmOKs) - getValueOrDefault(prev.ProcSnmp6.ReasmOKs)
	netStat.Ip6ReasmFails = getValueOrDefault(curr.ProcSnmp6.ReasmFails) - getValueOrDefault(prev.ProcSnmp6.ReasmFails)
	netStat.Ip6FragOKs = getValueOrDefault(curr.ProcSnmp6.FragOKs) - getValueOrDefault(prev.ProcSnmp6.FragOKs)
	netStat.Ip6FragFails = getValueOrDefault(curr.ProcSnmp6.FragFails) - getValueOrDefault(prev.ProcSnmp6.FragFails)
	netStat.Ip6FragCreates = getValueOrDefault(curr.ProcSnmp6.FragCreates) - getValueOrDefault(prev.ProcSnmp6.FragCreates)
	netStat.Ip6InMcastPkts = getValueOrDefault(curr.ProcSnmp6.InMcastPkts) - getValueOrDefault(prev.ProcSnmp6.InMcastPkts)
	netStat.Ip6OutMcastPkts = getValueOrDefault(curr.ProcSnmp6.OutMcastPkts) - getValueOrDefault(prev.ProcSnmp6.OutMcastPkts)
	netStat.Ip6InOctets = getValueOrDefault(curr.ProcSnmp6.InOctets) - getValueOrDefault(prev.ProcSnmp6.InOctets)
	netStat.Ip6OutOctets = getValueOrDefault(curr.ProcSnmp6.OutOctets) - getValueOrDefault(prev.ProcSnmp6.OutOctets)
	netStat.Ip6InMcastOctets = getValueOrDefault(curr.ProcSnmp6.InMcastOctets) - getValueOrDefault(prev.ProcSnmp6.InMcastOctets)
	netStat.Ip6OutMcastOctets = getValueOrDefault(curr.ProcSnmp6.OutMcastOctets) - getValueOrDefault(prev.ProcSnmp6.OutMcastOctets)
	netStat.Ip6InBcastOctets = getValueOrDefault(curr.ProcSnmp6.InBcastOctets) - getValueOrDefault(prev.ProcSnmp6.InBcastOctets)
	netStat.Ip6OutBcastOctets = getValueOrDefault(curr.ProcSnmp6.OutBcastOctets) - getValueOrDefault(prev.ProcSnmp6.OutBcastOctets)
	netStat.Ip6InNoECTPkts = getValueOrDefault(curr.ProcSnmp6.InNoECTPkts) - getValueOrDefault(prev.ProcSnmp6.InNoECTPkts)
	netStat.Ip6InECT1Pkts = getValueOrDefault(curr.ProcSnmp6.InECT1Pkts) - getValueOrDefault(prev.ProcSnmp6.InECT1Pkts)
	netStat.Ip6InECT0Pkts = getValueOrDefault(curr.ProcSnmp6.InECT0Pkts) - getValueOrDefault(prev.ProcSnmp6.InECT0Pkts)
	netStat.Ip6InCEPkts = getValueOrDefault(curr.ProcSnmp6.InCEPkts) - getValueOrDefault(prev.ProcSnmp6.InCEPkts)

	netStat.Icmp6InMsgs = getValueOrDefault(curr.ProcSnmp6.InMsgs) - getValueOrDefault(prev.ProcSnmp6.InMsgs)
	netStat.Icmp6InErrors = getValueOrDefault(curr.ProcSnmp6.Icmp6.InErrors) - getValueOrDefault(prev.ProcSnmp6.Icmp6.InErrors)
	netStat.Icmp6OutMsgs = getValueOrDefault(curr.ProcSnmp6.OutMsgs) - getValueOrDefault(prev.ProcSnmp6.OutMsgs)
	netStat.Icmp6OutErrors = getValueOrDefault(curr.ProcSnmp6.OutErrors) - getValueOrDefault(prev.ProcSnmp6.OutErrors)
	netStat.Icmp6InCsumErrors = getValueOrDefault(curr.ProcSnmp6.Icmp6.InCsumErrors) - getValueOrDefault(prev.ProcSnmp6.Icmp6.InCsumErrors)
	netStat.Icmp6InDestUnreachs = getValueOrDefault(curr.ProcSnmp6.InDestUnreachs) - getValueOrDefault(prev.ProcSnmp6.InDestUnreachs)
	netStat.Icmp6InPktTooBigs = getValueOrDefault(curr.ProcSnmp6.InPktTooBigs) - getValueOrDefault(prev.ProcSnmp6.InPktTooBigs)
	netStat.Icmp6InTimeExcds = getValueOrDefault(curr.ProcSnmp6.InTimeExcds) - getValueOrDefault(prev.ProcSnmp6.InTimeExcds)
	netStat.Icmp6InParmProblems = getValueOrDefault(curr.ProcSnmp6.InParmProblems) - getValueOrDefault(prev.ProcSnmp6.InParmProblems)
	netStat.Icmp6InEchos = getValueOrDefault(curr.ProcSnmp6.InEchos) - getValueOrDefault(prev.ProcSnmp6.InEchos)
	netStat.Icmp6InEchoReplies = getValueOrDefault(curr.ProcSnmp6.InEchoReplies) - getValueOrDefault(prev.ProcSnmp6.InEchoReplies)
	netStat.Icmp6InGroupMembQueries = getValueOrDefault(curr.ProcSnmp6.InGroupMembQueries) - getValueOrDefault(prev.ProcSnmp6.InGroupMembQueries)
	netStat.Icmp6InGroupMembResponses = getValueOrDefault(curr.ProcSnmp6.InGroupMembResponses) - getValueOrDefault(prev.ProcSnmp6.InGroupMembResponses)
	netStat.Icmp6InGroupMembReductions = getValueOrDefault(curr.ProcSnmp6.InGroupMembReductions) - getValueOrDefault(prev.ProcSnmp6.InGroupMembReductions)
	netStat.Icmp6InRouterSolicits = getValueOrDefault(curr.ProcSnmp6.InRouterSolicits) - getValueOrDefault(prev.ProcSnmp6.InRouterSolicits)
	netStat.Icmp6InRouterAdvertisements = getValueOrDefault(curr.ProcSnmp6.InRouterAdvertisements) - getValueOrDefault(prev.ProcSnmp6.InRouterAdvertisements)
	netStat.Icmp6InNeighborSolicits = getValueOrDefault(curr.ProcSnmp6.InNeighborSolicits) - getValueOrDefault(prev.ProcSnmp6.InNeighborSolicits)
	netStat.Icmp6InNeighborAdvertisements = getValueOrDefault(curr.ProcSnmp6.InNeighborAdvertisements) - getValueOrDefault(prev.ProcSnmp6.InNeighborAdvertisements)
	netStat.Icmp6InRedirects = getValueOrDefault(curr.ProcSnmp6.InRedirects) - getValueOrDefault(prev.ProcSnmp6.InRedirects)
	netStat.Icmp6InMLDv2Reports = getValueOrDefault(curr.ProcSnmp6.InMLDv2Reports) - getValueOrDefault(prev.ProcSnmp6.InMLDv2Reports)
	netStat.Icmp6OutDestUnreachs = getValueOrDefault(curr.ProcSnmp6.OutDestUnreachs) - getValueOrDefault(prev.ProcSnmp6.OutDestUnreachs)
	netStat.Icmp6OutPktTooBigs = getValueOrDefault(curr.ProcSnmp6.OutPktTooBigs) - getValueOrDefault(prev.ProcSnmp6.OutPktTooBigs)
	netStat.Icmp6OutTimeExcds = getValueOrDefault(curr.ProcSnmp6.OutTimeExcds) - getValueOrDefault(prev.ProcSnmp6.OutTimeExcds)
	netStat.Icmp6OutParmProblems = getValueOrDefault(curr.ProcSnmp6.OutParmProblems) - getValueOrDefault(prev.ProcSnmp6.OutParmProblems)
	netStat.Icmp6OutEchos = getValueOrDefault(curr.ProcSnmp6.OutEchos) - getValueOrDefault(prev.ProcSnmp6.OutEchos)
	netStat.Icmp6OutEchoReplies = getValueOrDefault(curr.ProcSnmp6.OutEchoReplies) - getValueOrDefault(prev.ProcSnmp6.OutEchoReplies)
	netStat.Icmp6OutGroupMembQueries = getValueOrDefault(curr.ProcSnmp6.OutGroupMembQueries) - getValueOrDefault(prev.ProcSnmp6.OutGroupMembQueries)
	netStat.Icmp6OutGroupMembResponses = getValueOrDefault(curr.ProcSnmp6.OutGroupMembResponses) - getValueOrDefault(prev.ProcSnmp6.OutGroupMembResponses)
	netStat.Icmp6OutGroupMembReductions = getValueOrDefault(curr.ProcSnmp6.OutGroupMembReductions) - getValueOrDefault(prev.ProcSnmp6.OutGroupMembReductions)
	netStat.Icmp6OutRouterSolicits = getValueOrDefault(curr.ProcSnmp6.OutRouterSolicits) - getValueOrDefault(prev.ProcSnmp6.OutRouterSolicits)
	netStat.Icmp6OutRouterAdvertisements = getValueOrDefault(curr.ProcSnmp6.OutRouterAdvertisements) - getValueOrDefault(prev.ProcSnmp6.OutRouterAdvertisements)
	netStat.Icmp6OutNeighborSolicits = getValueOrDefault(curr.ProcSnmp6.OutNeighborSolicits) - getValueOrDefault(prev.ProcSnmp6.OutNeighborSolicits)
	netStat.Icmp6OutNeighborAdvertisements = getValueOrDefault(curr.ProcSnmp6.OutNeighborAdvertisements) - getValueOrDefault(prev.ProcSnmp6.OutNeighborAdvertisements)
	netStat.Icmp6OutRedirects = getValueOrDefault(curr.ProcSnmp6.OutRedirects) - getValueOrDefault(prev.ProcSnmp6.OutRedirects)
	netStat.Icmp6OutMLDv2Reports = getValueOrDefault(curr.ProcSnmp6.OutMLDv2Reports) - getValueOrDefault(prev.ProcSnmp6.OutMLDv2Reports)
	netStat.Icmp6InType1 = getValueOrDefault(curr.ProcSnmp6.InType1) - getValueOrDefault(prev.ProcSnmp6.InType1)
	netStat.Icmp6InType134 = getValueOrDefault(curr.ProcSnmp6.InType134) - getValueOrDefault(prev.ProcSnmp6.InType134)
	netStat.Icmp6InType135 = getValueOrDefault(curr.ProcSnmp6.InType135) - getValueOrDefault(prev.ProcSnmp6.InType135)
	netStat.Icmp6InType136 = getValueOrDefault(curr.ProcSnmp6.InType136) - getValueOrDefault(prev.ProcSnmp6.InType136)
	netStat.Icmp6InType143 = getValueOrDefault(curr.ProcSnmp6.InType143) - getValueOrDefault(prev.ProcSnmp6.InType143)
	netStat.Icmp6OutType133 = getValueOrDefault(curr.ProcSnmp6.OutType133) - getValueOrDefault(prev.ProcSnmp6.OutType133)
	netStat.Icmp6OutType135 = getValueOrDefault(curr.ProcSnmp6.OutType135) - getValueOrDefault(prev.ProcSnmp6.OutType135)
	netStat.Icmp6OutType136 = getValueOrDefault(curr.ProcSnmp6.OutType136) - getValueOrDefault(prev.ProcSnmp6.OutType136)
	netStat.Icmp6OutType143 = getValueOrDefault(curr.ProcSnmp6.OutType143) - getValueOrDefault(prev.ProcSnmp6.OutType143)

	netStat.Udp6InDatagrams = getValueOrDefault(curr.ProcSnmp6.Udp6.InDatagrams) - getValueOrDefault(prev.ProcSnmp6.Udp6.InDatagrams)
	netStat.Udp6NoPorts = getValueOrDefault(curr.ProcSnmp6.Udp6.NoPorts) - getValueOrDefault(prev.ProcSnmp6.Udp6.NoPorts)
	netStat.Udp6InErrors = getValueOrDefault(curr.ProcSnmp6.Udp6.InErrors) - getValueOrDefault(prev.ProcSnmp6.Udp6.InErrors)
	netStat.Udp6OutDatagrams = getValueOrDefault(curr.ProcSnmp6.Udp6.OutDatagrams) - getValueOrDefault(prev.ProcSnmp6.Udp6.OutDatagrams)
	netStat.Udp6RcvbufErrors = getValueOrDefault(curr.ProcSnmp6.Udp6.RcvbufErrors) - getValueOrDefault(prev.ProcSnmp6.Udp6.RcvbufErrors)
	netStat.Udp6SndbufErrors = getValueOrDefault(curr.ProcSnmp6.Udp6.SndbufErrors) - getValueOrDefault(prev.ProcSnmp6.Udp6.SndbufErrors)
	netStat.Udp6InCsumErrors = getValueOrDefault(curr.ProcSnmp6.Udp6.InCsumErrors) - getValueOrDefault(prev.ProcSnmp6.Udp6.InCsumErrors)
	netStat.Udp6IgnoredMulti = getValueOrDefault(curr.ProcSnmp6.Udp6.IgnoredMulti) - getValueOrDefault(prev.ProcSnmp6.Udp6.IgnoredMulti)

	netStat.UdpLite6InDatagrams = getValueOrDefault(curr.ProcSnmp6.UdpLite6.InDatagrams) - getValueOrDefault(prev.ProcSnmp6.UdpLite6.InDatagrams)
	netStat.UdpLite6NoPorts = getValueOrDefault(curr.ProcSnmp6.UdpLite6.NoPorts) - getValueOrDefault(prev.ProcSnmp6.UdpLite6.NoPorts)
	netStat.UdpLite6InErrors = getValueOrDefault(curr.ProcSnmp6.UdpLite6.InErrors) - getValueOrDefault(prev.ProcSnmp6.UdpLite6.InErrors)
	netStat.UdpLite6OutDatagrams = getValueOrDefault(curr.ProcSnmp6.UdpLite6.OutDatagrams) - getValueOrDefault(prev.ProcSnmp6.UdpLite6.OutDatagrams)
	netStat.UdpLite6RcvbufErrors = getValueOrDefault(curr.ProcSnmp6.UdpLite6.RcvbufErrors) - getValueOrDefault(prev.ProcSnmp6.UdpLite6.RcvbufErrors)
	netStat.UdpLite6SndbufErrors = getValueOrDefault(curr.ProcSnmp6.UdpLite6.SndbufErrors) - getValueOrDefault(prev.ProcSnmp6.UdpLite6.SndbufErrors)
	netStat.UdpLite6InCsumErrors = getValueOrDefault(curr.ProcSnmp6.UdpLite6.InCsumErrors) - getValueOrDefault(prev.ProcSnmp6.UdpLite6.InCsumErrors)

	netStat.TcpExtSyncookiesSent = getValueOrDefault(curr.ProcNetstat.SyncookiesSent) - getValueOrDefault(prev.ProcNetstat.SyncookiesSent)
	netStat.TcpExtSyncookiesRecv = getValueOrDefault(curr.ProcNetstat.SyncookiesRecv) - getValueOrDefault(prev.ProcNetstat.SyncookiesRecv)
	netStat.TcpExtSyncookiesFailed = getValueOrDefault(curr.ProcNetstat.SyncookiesFailed) - getValueOrDefault(prev.ProcNetstat.SyncookiesFailed)
	netStat.TcpExtEmbryonicRsts = getValueOrDefault(curr.ProcNetstat.EmbryonicRsts) - getValueOrDefault(prev.ProcNetstat.EmbryonicRsts)
	netStat.TcpExtPruneCalled = getValueOrDefault(curr.ProcNetstat.PruneCalled) - getValueOrDefault(prev.ProcNetstat.PruneCalled)
	netStat.TcpExtRcvPruned = getValueOrDefault(curr.ProcNetstat.RcvPruned) - getValueOrDefault(prev.ProcNetstat.RcvPruned)
	netStat.TcpExtOfoPruned = getValueOrDefault(curr.ProcNetstat.OfoPruned) - getValueOrDefault(prev.ProcNetstat.OfoPruned)
	netStat.TcpExtOutOfWindowIcmps = getValueOrDefault(curr.ProcNetstat.OutOfWindowIcmps) - getValueOrDefault(prev.ProcNetstat.OutOfWindowIcmps)
	netStat.TcpExtLockDroppedIcmps = getValueOrDefault(curr.ProcNetstat.LockDroppedIcmps) - getValueOrDefault(prev.ProcNetstat.LockDroppedIcmps)
	netStat.TcpExtArpFilter = getValueOrDefault(curr.ProcNetstat.ArpFilter) - getValueOrDefault(prev.ProcNetstat.ArpFilter)
	netStat.TcpExtTW = getValueOrDefault(curr.ProcNetstat.TW) - getValueOrDefault(prev.ProcNetstat.TW)
	netStat.TcpExtTWRecycled = getValueOrDefault(curr.ProcNetstat.TWRecycled) - getValueOrDefault(prev.ProcNetstat.TWRecycled)
	netStat.TcpExtTWKilled = getValueOrDefault(curr.ProcNetstat.TWKilled) - getValueOrDefault(prev.ProcNetstat.TWKilled)
	netStat.TcpExtPAWSActive = getValueOrDefault(curr.ProcNetstat.PAWSActive) - getValueOrDefault(prev.ProcNetstat.PAWSActive)
	netStat.TcpExtPAWSEstab = getValueOrDefault(curr.ProcNetstat.PAWSEstab) - getValueOrDefault(prev.ProcNetstat.PAWSEstab)
	netStat.TcpExtDelayedACKs = getValueOrDefault(curr.ProcNetstat.DelayedACKs) - getValueOrDefault(prev.ProcNetstat.DelayedACKs)
	netStat.TcpExtDelayedACKLocked = getValueOrDefault(curr.ProcNetstat.DelayedACKLocked) - getValueOrDefault(prev.ProcNetstat.DelayedACKLocked)
	netStat.TcpExtDelayedACKLost = getValueOrDefault(curr.ProcNetstat.DelayedACKLost) - getValueOrDefault(prev.ProcNetstat.DelayedACKLost)
	netStat.TcpExtListenOverflows = getValueOrDefault(curr.ProcNetstat.ListenOverflows) - getValueOrDefault(prev.ProcNetstat.ListenOverflows)
	netStat.TcpExtListenDrops = getValueOrDefault(curr.ProcNetstat.ListenDrops) - getValueOrDefault(prev.ProcNetstat.ListenDrops)
	netStat.TcpExtTCPHPHits = getValueOrDefault(curr.ProcNetstat.TCPHPHits) - getValueOrDefault(prev.ProcNetstat.TCPHPHits)
	netStat.TcpExtTCPPureAcks = getValueOrDefault(curr.ProcNetstat.TCPPureAcks) - getValueOrDefault(prev.ProcNetstat.TCPPureAcks)
	netStat.TcpExtTCPHPAcks = getValueOrDefault(curr.ProcNetstat.TCPHPAcks) - getValueOrDefault(prev.ProcNetstat.TCPHPAcks)
	netStat.TcpExtTCPRenoRecovery = getValueOrDefault(curr.ProcNetstat.TCPRenoRecovery) - getValueOrDefault(prev.ProcNetstat.TCPRenoRecovery)
	netStat.TcpExtTCPSackRecovery = getValueOrDefault(curr.ProcNetstat.TCPSackRecovery) - getValueOrDefault(prev.ProcNetstat.TCPSackRecovery)
	netStat.TcpExtTCPSACKReneging = getValueOrDefault(curr.ProcNetstat.TCPSACKReneging) - getValueOrDefault(prev.ProcNetstat.TCPSACKReneging)
	netStat.TcpExtTCPSACKReorder = getValueOrDefault(curr.ProcNetstat.TCPSACKReorder) - getValueOrDefault(prev.ProcNetstat.TCPSACKReorder)
	netStat.TcpExtTCPRenoReorder = getValueOrDefault(curr.ProcNetstat.TCPRenoReorder) - getValueOrDefault(prev.ProcNetstat.TCPRenoReorder)
	netStat.TcpExtTCPTSReorder = getValueOrDefault(curr.ProcNetstat.TCPTSReorder) - getValueOrDefault(prev.ProcNetstat.TCPTSReorder)
	netStat.TcpExtTCPFullUndo = getValueOrDefault(curr.ProcNetstat.TCPFullUndo) - getValueOrDefault(prev.ProcNetstat.TCPFullUndo)
	netStat.TcpExtTCPPartialUndo = getValueOrDefault(curr.ProcNetstat.TCPPartialUndo) - getValueOrDefault(prev.ProcNetstat.TCPPartialUndo)
	netStat.TcpExtTCPDSACKUndo = getValueOrDefault(curr.ProcNetstat.TCPDSACKUndo) - getValueOrDefault(prev.ProcNetstat.TCPDSACKUndo)
	netStat.TcpExtTCPLossUndo = getValueOrDefault(curr.ProcNetstat.TCPLossUndo) - getValueOrDefault(prev.ProcNetstat.TCPLossUndo)
	netStat.TcpExtTCPLostRetransmit = getValueOrDefault(curr.ProcNetstat.TCPLostRetransmit) - getValueOrDefault(prev.ProcNetstat.TCPLostRetransmit)
	netStat.TcpExtTCPRenoFailures = getValueOrDefault(curr.ProcNetstat.TCPRenoFailures) - getValueOrDefault(prev.ProcNetstat.TCPRenoFailures)
	netStat.TcpExtTCPSackFailures = getValueOrDefault(curr.ProcNetstat.TCPSackFailures) - getValueOrDefault(prev.ProcNetstat.TCPSackFailures)
	netStat.TcpExtTCPLossFailures = getValueOrDefault(curr.ProcNetstat.TCPLossFailures) - getValueOrDefault(prev.ProcNetstat.TCPLossFailures)
	netStat.TcpExtTCPFastRetrans = getValueOrDefault(curr.ProcNetstat.TCPFastRetrans) - getValueOrDefault(prev.ProcNetstat.TCPFastRetrans)
	netStat.TcpExtTCPSlowStartRetrans = getValueOrDefault(curr.ProcNetstat.TCPSlowStartRetrans) - getValueOrDefault(prev.ProcNetstat.TCPSlowStartRetrans)
	netStat.TcpExtTCPTimeouts = getValueOrDefault(curr.ProcNetstat.TCPTimeouts) - getValueOrDefault(prev.ProcNetstat.TCPTimeouts)
	netStat.TcpExtTCPLossProbes = getValueOrDefault(curr.ProcNetstat.TCPLossProbes) - getValueOrDefault(prev.ProcNetstat.TCPLossProbes)
	netStat.TcpExtTCPLossProbeRecovery = getValueOrDefault(curr.ProcNetstat.TCPLossProbeRecovery) - getValueOrDefault(prev.ProcNetstat.TCPLossProbeRecovery)
	netStat.TcpExtTCPRenoRecoveryFail = getValueOrDefault(curr.ProcNetstat.TCPRenoRecoveryFail) - getValueOrDefault(prev.ProcNetstat.TCPRenoRecoveryFail)
	netStat.TcpExtTCPSackRecoveryFail = getValueOrDefault(curr.ProcNetstat.TCPSackRecoveryFail) - getValueOrDefault(prev.ProcNetstat.TCPSackRecoveryFail)
	netStat.TcpExtTCPRcvCollapsed = getValueOrDefault(curr.ProcNetstat.TCPRcvCollapsed) - getValueOrDefault(prev.ProcNetstat.TCPRcvCollapsed)
	netStat.TcpExtTCPDSACKOldSent = getValueOrDefault(curr.ProcNetstat.TCPDSACKOldSent) - getValueOrDefault(prev.ProcNetstat.TCPDSACKOldSent)
	netStat.TcpExtTCPDSACKOfoSent = getValueOrDefault(curr.ProcNetstat.TCPDSACKOfoSent) - getValueOrDefault(prev.ProcNetstat.TCPDSACKOfoSent)
	netStat.TcpExtTCPDSACKRecv = getValueOrDefault(curr.ProcNetstat.TCPDSACKRecv) - getValueOrDefault(prev.ProcNetstat.TCPDSACKRecv)
	netStat.TcpExtTCPDSACKOfoRecv = getValueOrDefault(curr.ProcNetstat.TCPDSACKOfoRecv) - getValueOrDefault(prev.ProcNetstat.TCPDSACKOfoRecv)
	netStat.TcpExtTCPAbortOnData = getValueOrDefault(curr.ProcNetstat.TCPAbortOnData) - getValueOrDefault(prev.ProcNetstat.TCPAbortOnData)
	netStat.TcpExtTCPAbortOnClose = getValueOrDefault(curr.ProcNetstat.TCPAbortOnClose) - getValueOrDefault(prev.ProcNetstat.TCPAbortOnClose)
	netStat.TcpExtTCPAbortOnMemory = getValueOrDefault(curr.ProcNetstat.TCPAbortOnMemory) - getValueOrDefault(prev.ProcNetstat.TCPAbortOnMemory)
	netStat.TcpExtTCPAbortOnTimeout = getValueOrDefault(curr.ProcNetstat.TCPAbortOnTimeout) - getValueOrDefault(prev.ProcNetstat.TCPAbortOnTimeout)
	netStat.TcpExtTCPAbortOnLinger = getValueOrDefault(curr.ProcNetstat.TCPAbortOnLinger) - getValueOrDefault(prev.ProcNetstat.TCPAbortOnLinger)
	netStat.TcpExtTCPAbortFailed = getValueOrDefault(curr.ProcNetstat.TCPAbortFailed) - getValueOrDefault(prev.ProcNetstat.TCPAbortFailed)
	netStat.TcpExtTCPMemoryPressures = getValueOrDefault(curr.ProcNetstat.TCPMemoryPressures) - getValueOrDefault(prev.ProcNetstat.TCPMemoryPressures)
	netStat.TcpExtTCPMemoryPressuresChrono = getValueOrDefault(curr.ProcNetstat.TCPMemoryPressuresChrono) - getValueOrDefault(prev.ProcNetstat.TCPMemoryPressuresChrono)
	netStat.TcpExtTCPSACKDiscard = getValueOrDefault(curr.ProcNetstat.TCPSACKDiscard) - getValueOrDefault(prev.ProcNetstat.TCPSACKDiscard)
	netStat.TcpExtTCPDSACKIgnoredOld = getValueOrDefault(curr.ProcNetstat.TCPDSACKIgnoredOld) - getValueOrDefault(prev.ProcNetstat.TCPDSACKIgnoredOld)
	netStat.TcpExtTCPDSACKIgnoredNoUndo = getValueOrDefault(curr.ProcNetstat.TCPDSACKIgnoredNoUndo) - getValueOrDefault(prev.ProcNetstat.TCPDSACKIgnoredNoUndo)
	netStat.TcpExtTCPSpuriousRTOs = getValueOrDefault(curr.ProcNetstat.TCPSpuriousRTOs) - getValueOrDefault(prev.ProcNetstat.TCPSpuriousRTOs)
	netStat.TcpExtTCPMD5NotFound = getValueOrDefault(curr.ProcNetstat.TCPMD5NotFound) - getValueOrDefault(prev.ProcNetstat.TCPMD5NotFound)
	netStat.TcpExtTCPMD5Unexpected = getValueOrDefault(curr.ProcNetstat.TCPMD5Unexpected) - getValueOrDefault(prev.ProcNetstat.TCPMD5Unexpected)
	netStat.TcpExtTCPMD5Failure = getValueOrDefault(curr.ProcNetstat.TCPMD5Failure) - getValueOrDefault(prev.ProcNetstat.TCPMD5Failure)
	netStat.TcpExtTCPSackShifted = getValueOrDefault(curr.ProcNetstat.TCPSackShifted) - getValueOrDefault(prev.ProcNetstat.TCPSackShifted)
	netStat.TcpExtTCPSackMerged = getValueOrDefault(curr.ProcNetstat.TCPSackMerged) - getValueOrDefault(prev.ProcNetstat.TCPSackMerged)
	netStat.TcpExtTCPSackShiftFallback = getValueOrDefault(curr.ProcNetstat.TCPSackShiftFallback) - getValueOrDefault(prev.ProcNetstat.TCPSackShiftFallback)
	netStat.TcpExtTCPBacklogDrop = getValueOrDefault(curr.ProcNetstat.TCPBacklogDrop) - getValueOrDefault(prev.ProcNetstat.TCPBacklogDrop)
	netStat.TcpExtPFMemallocDrop = getValueOrDefault(curr.ProcNetstat.PFMemallocDrop) - getValueOrDefault(prev.ProcNetstat.PFMemallocDrop)
	netStat.TcpExtTCPMinTTLDrop = getValueOrDefault(curr.ProcNetstat.TCPMinTTLDrop) - getValueOrDefault(prev.ProcNetstat.TCPMinTTLDrop)
	netStat.TcpExtTCPDeferAcceptDrop = getValueOrDefault(curr.ProcNetstat.TCPDeferAcceptDrop) - getValueOrDefault(prev.ProcNetstat.TCPDeferAcceptDrop)
	netStat.TcpExtIPReversePathFilter = getValueOrDefault(curr.ProcNetstat.IPReversePathFilter) - getValueOrDefault(prev.ProcNetstat.IPReversePathFilter)
	netStat.TcpExtTCPTimeWaitOverflow = getValueOrDefault(curr.ProcNetstat.TCPTimeWaitOverflow) - getValueOrDefault(prev.ProcNetstat.TCPTimeWaitOverflow)
	netStat.TcpExtTCPReqQFullDoCookies = getValueOrDefault(curr.ProcNetstat.TCPReqQFullDoCookies) - getValueOrDefault(prev.ProcNetstat.TCPReqQFullDoCookies)
	netStat.TcpExtTCPReqQFullDrop = getValueOrDefault(curr.ProcNetstat.TCPReqQFullDrop) - getValueOrDefault(prev.ProcNetstat.TCPReqQFullDrop)
	netStat.TcpExtTCPRetransFail = getValueOrDefault(curr.ProcNetstat.TCPRetransFail) - getValueOrDefault(prev.ProcNetstat.TCPRetransFail)
	netStat.TcpExtTCPRcvCoalesce = getValueOrDefault(curr.ProcNetstat.TCPRcvCoalesce) - getValueOrDefault(prev.ProcNetstat.TCPRcvCoalesce)
	netStat.TcpExtTCPRcvQDrop = getValueOrDefault(curr.ProcNetstat.TCPRcvQDrop) - getValueOrDefault(prev.ProcNetstat.TCPRcvQDrop)
	netStat.TcpExtTCPOFOQueue = getValueOrDefault(curr.ProcNetstat.TCPOFOQueue) - getValueOrDefault(prev.ProcNetstat.TCPOFOQueue)
	netStat.TcpExtTCPOFODrop = getValueOrDefault(curr.ProcNetstat.TCPOFODrop) - getValueOrDefault(prev.ProcNetstat.TCPOFODrop)
	netStat.TcpExtTCPOFOMerge = getValueOrDefault(curr.ProcNetstat.TCPOFOMerge) - getValueOrDefault(prev.ProcNetstat.TCPOFOMerge)
	netStat.TcpExtTCPChallengeACK = getValueOrDefault(curr.ProcNetstat.TCPChallengeACK) - getValueOrDefault(prev.ProcNetstat.TCPChallengeACK)
	netStat.TcpExtTCPSYNChallenge = getValueOrDefault(curr.ProcNetstat.TCPSYNChallenge) - getValueOrDefault(prev.ProcNetstat.TCPSYNChallenge)
	netStat.TcpExtTCPFastOpenActive = getValueOrDefault(curr.ProcNetstat.TCPFastOpenActive) - getValueOrDefault(prev.ProcNetstat.TCPFastOpenActive)
	netStat.TcpExtTCPFastOpenActiveFail = getValueOrDefault(curr.ProcNetstat.TCPFastOpenActiveFail) - getValueOrDefault(prev.ProcNetstat.TCPFastOpenActiveFail)
	netStat.TcpExtTCPFastOpenPassive = getValueOrDefault(curr.ProcNetstat.TCPFastOpenPassive) - getValueOrDefault(prev.ProcNetstat.TCPFastOpenPassive)
	netStat.TcpExtTCPFastOpenPassiveFail = getValueOrDefault(curr.ProcNetstat.TCPFastOpenPassiveFail) - getValueOrDefault(prev.ProcNetstat.TCPFastOpenPassiveFail)
	netStat.TcpExtTCPFastOpenListenOverflow = getValueOrDefault(curr.ProcNetstat.TCPFastOpenListenOverflow) - getValueOrDefault(prev.ProcNetstat.TCPFastOpenListenOverflow)
	netStat.TcpExtTCPFastOpenCookieReqd = getValueOrDefault(curr.ProcNetstat.TCPFastOpenCookieReqd) - getValueOrDefault(prev.ProcNetstat.TCPFastOpenCookieReqd)
	netStat.TcpExtTCPFastOpenBlackhole = getValueOrDefault(curr.ProcNetstat.TCPFastOpenBlackhole) - getValueOrDefault(prev.ProcNetstat.TCPFastOpenBlackhole)
	netStat.TcpExtTCPSpuriousRtxHostQueues = getValueOrDefault(curr.ProcNetstat.TCPSpuriousRtxHostQueues) - getValueOrDefault(prev.ProcNetstat.TCPSpuriousRtxHostQueues)
	netStat.TcpExtBusyPollRxPackets = getValueOrDefault(curr.ProcNetstat.BusyPollRxPackets) - getValueOrDefault(prev.ProcNetstat.BusyPollRxPackets)
	netStat.TcpExtTCPAutoCorking = getValueOrDefault(curr.ProcNetstat.TCPAutoCorking) - getValueOrDefault(prev.ProcNetstat.TCPAutoCorking)
	netStat.TcpExtTCPFromZeroWindowAdv = getValueOrDefault(curr.ProcNetstat.TCPFromZeroWindowAdv) - getValueOrDefault(prev.ProcNetstat.TCPFromZeroWindowAdv)
	netStat.TcpExtTCPToZeroWindowAdv = getValueOrDefault(curr.ProcNetstat.TCPToZeroWindowAdv) - getValueOrDefault(prev.ProcNetstat.TCPToZeroWindowAdv)
	netStat.TcpExtTCPWantZeroWindowAdv = getValueOrDefault(curr.ProcNetstat.TCPWantZeroWindowAdv) - getValueOrDefault(prev.ProcNetstat.TCPWantZeroWindowAdv)
	netStat.TcpExtTCPSynRetrans = getValueOrDefault(curr.ProcNetstat.TCPSynRetrans) - getValueOrDefault(prev.ProcNetstat.TCPSynRetrans)
	netStat.TcpExtTCPOrigDataSent = getValueOrDefault(curr.ProcNetstat.TCPOrigDataSent) - getValueOrDefault(prev.ProcNetstat.TCPOrigDataSent)
	netStat.TcpExtTCPHystartTrainDetect = getValueOrDefault(curr.ProcNetstat.TCPHystartTrainDetect) - getValueOrDefault(prev.ProcNetstat.TCPHystartTrainDetect)
	netStat.TcpExtTCPHystartTrainCwnd = getValueOrDefault(curr.ProcNetstat.TCPHystartTrainCwnd) - getValueOrDefault(prev.ProcNetstat.TCPHystartTrainCwnd)
	netStat.TcpExtTCPHystartDelayDetect = getValueOrDefault(curr.ProcNetstat.TCPHystartDelayDetect) - getValueOrDefault(prev.ProcNetstat.TCPHystartDelayDetect)
	netStat.TcpExtTCPHystartDelayCwnd = getValueOrDefault(curr.ProcNetstat.TCPHystartDelayCwnd) - getValueOrDefault(prev.ProcNetstat.TCPHystartDelayCwnd)
	netStat.TcpExtTCPACKSkippedSynRecv = getValueOrDefault(curr.ProcNetstat.TCPACKSkippedSynRecv) - getValueOrDefault(prev.ProcNetstat.TCPACKSkippedSynRecv)
	netStat.TcpExtTCPACKSkippedPAWS = getValueOrDefault(curr.ProcNetstat.TCPACKSkippedPAWS) - getValueOrDefault(prev.ProcNetstat.TCPACKSkippedPAWS)
	netStat.TcpExtTCPACKSkippedSeq = getValueOrDefault(curr.ProcNetstat.TCPACKSkippedSeq) - getValueOrDefault(prev.ProcNetstat.TCPACKSkippedSeq)
	netStat.TcpExtTCPACKSkippedFinWait2 = getValueOrDefault(curr.ProcNetstat.TCPACKSkippedFinWait2) - getValueOrDefault(prev.ProcNetstat.TCPACKSkippedFinWait2)
	netStat.TcpExtTCPACKSkippedTimeWait = getValueOrDefault(curr.ProcNetstat.TCPACKSkippedTimeWait) - getValueOrDefault(prev.ProcNetstat.TCPACKSkippedTimeWait)
	netStat.TcpExtTCPACKSkippedChallenge = getValueOrDefault(curr.ProcNetstat.TCPACKSkippedChallenge) - getValueOrDefault(prev.ProcNetstat.TCPACKSkippedChallenge)
	netStat.TcpExtTCPWinProbe = getValueOrDefault(curr.ProcNetstat.TCPWinProbe) - getValueOrDefault(prev.ProcNetstat.TCPWinProbe)
	netStat.TcpExtTCPKeepAlive = getValueOrDefault(curr.ProcNetstat.TCPKeepAlive) - getValueOrDefault(prev.ProcNetstat.TCPKeepAlive)
	netStat.TcpExtTCPMTUPFail = getValueOrDefault(curr.ProcNetstat.TCPMTUPFail) - getValueOrDefault(prev.ProcNetstat.TCPMTUPFail)
	netStat.TcpExtTCPMTUPSuccess = getValueOrDefault(curr.ProcNetstat.TCPMTUPSuccess) - getValueOrDefault(prev.ProcNetstat.TCPMTUPSuccess)
	netStat.TcpExtTCPWqueueTooBig = getValueOrDefault(curr.ProcNetstat.TCPWqueueTooBig) - getValueOrDefault(prev.ProcNetstat.TCPWqueueTooBig)

	netStat.IpExtInNoRoutes = getValueOrDefault(curr.ProcNetstat.InNoRoutes) - getValueOrDefault(prev.ProcNetstat.InNoRoutes)
	netStat.IpExtInTruncatedPkts = getValueOrDefault(curr.ProcNetstat.InTruncatedPkts) - getValueOrDefault(prev.ProcNetstat.InTruncatedPkts)
	netStat.IpExtInMcastPkts = getValueOrDefault(curr.ProcNetstat.InMcastPkts) - getValueOrDefault(prev.ProcNetstat.InMcastPkts)
	netStat.IpExtOutMcastPkts = getValueOrDefault(curr.ProcNetstat.OutMcastPkts) - getValueOrDefault(prev.ProcNetstat.OutMcastPkts)
	netStat.IpExtInBcastPkts = getValueOrDefault(curr.ProcNetstat.InBcastPkts) - getValueOrDefault(prev.ProcNetstat.InBcastPkts)
	netStat.IpExtOutBcastPkts = getValueOrDefault(curr.ProcNetstat.OutBcastPkts) - getValueOrDefault(prev.ProcNetstat.OutBcastPkts)
	netStat.IpExtInOctets = getValueOrDefault(curr.ProcNetstat.InOctets) - getValueOrDefault(prev.ProcNetstat.InOctets)
	netStat.IpExtOutOctets = getValueOrDefault(curr.ProcNetstat.OutOctets) - getValueOrDefault(prev.ProcNetstat.OutOctets)
	netStat.IpExtInMcastOctets = getValueOrDefault(curr.ProcNetstat.InMcastOctets) - getValueOrDefault(prev.ProcNetstat.InMcastOctets)
	netStat.IpExtOutMcastOctets = getValueOrDefault(curr.ProcNetstat.OutMcastOctets) - getValueOrDefault(prev.ProcNetstat.OutMcastOctets)
	netStat.IpExtInBcastOctets = getValueOrDefault(curr.ProcNetstat.InBcastOctets) - getValueOrDefault(prev.ProcNetstat.InBcastOctets)
	netStat.IpExtOutBcastOctets = getValueOrDefault(curr.ProcNetstat.OutBcastOctets) - getValueOrDefault(prev.ProcNetstat.OutBcastOctets)
	netStat.IpExtInCsumErrors = getValueOrDefault(curr.ProcNetstat.InCsumErrors) - getValueOrDefault(prev.ProcNetstat.InCsumErrors)
	netStat.IpExtInNoECTPkts = getValueOrDefault(curr.ProcNetstat.InNoECTPkts) - getValueOrDefault(prev.ProcNetstat.InNoECTPkts)
	netStat.IpExtInECT1Pkts = getValueOrDefault(curr.ProcNetstat.InECT1Pkts) - getValueOrDefault(prev.ProcNetstat.InECT1Pkts)
	netStat.IpExtInECT0Pkts = getValueOrDefault(curr.ProcNetstat.InECT0Pkts) - getValueOrDefault(prev.ProcNetstat.InECT0Pkts)
	netStat.IpExtInCEPkts = getValueOrDefault(curr.ProcNetstat.InCEPkts) - getValueOrDefault(prev.ProcNetstat.InCEPkts)
	netStat.IpExtReasmOverlaps = getValueOrDefault(curr.ProcNetstat.ReasmOverlaps) - getValueOrDefault(prev.ProcNetstat.ReasmOverlaps)

}

func (netStat *NetStat) Dump(timeStamp int64, config RenderConfig, opt DumpOption) {

	dateTime := time.Unix(timeStamp, 0).Format(time.RFC3339)
	switch opt.Format {
	case "text":
		config.SetFixWidth(true)

		row := strings.Builder{}
		row.WriteString(dateTime)
		for _, f := range opt.Fields {
			renderValue := netStat.GetRenderValue(config, f)
			if f == opt.SelectField && opt.Filter != nil {
				if opt.Filter.MatchString(renderValue) == false {
					continue
				}
			}
			row.WriteString(" ")
			row.WriteString(renderValue)
		}
		row.WriteString("\n")

		opt.Output.WriteString(row.String())

	case "json":
		t := []any{}
		row := make(map[string]string)
		row["Timestamp"] = dateTime
		for _, f := range opt.Fields {
			renderValue := netStat.GetRenderValue(config, f)
			if f == opt.SelectField && opt.Filter != nil {
				if opt.Filter.MatchString(renderValue) == false {
					continue
				}
			}
			row[config[f].Name] = renderValue
		}
		t = append(t, row)

		b, _ := json.Marshal(t)
		opt.Output.Write(b)
	}
}
