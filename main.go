package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"sync/atomic"
	"time"
)

var (
	ipv      int
	n        int
	coloinfo = `{"ADL":{"CN":"阿德莱德,南澳州,澳大利亚","EN":"Adelaide, SA, Australia"},"AKL":{"CN":"奥克兰,新西兰","EN":"Auckland, New Zealand"},"ALG":{"CN":"阿尔及尔,阿尔及利亚","EN":"Algiers, Algeria"},"AMM":{"CN":"安曼,约旦","EN":"Amman, Jordan"},"AMS":{"CN":"阿姆斯特丹,荷兰","EN":"Amsterdam, Netherlands"},"ARI":{"CN":"阿里卡,智利","EN":"Arica, Chile"},"ARN":{"CN":"斯德哥尔摩,瑞典","EN":"Stockholm, Sweden"},"ASU":{"CN":"亚松森,巴拉圭","EN":"Asunción, Paraguay"},"ATH":{"CN":"雅典,希腊","EN":"Athens, Greece"},"ATL":{"CN":"亚特兰大,佐治亚州,美国","EN":"Atlanta, GA, United States"},"BAH":{"CN":"麦纳麦,巴林","EN":"Manama, Bahrain"},"BCN":{"CN":"巴塞罗那,西班牙","EN":"Barcelona, Spain"},"BEG":{"CN":"贝尔格莱德,塞尔维亚","EN":"Belgrade, Serbia"},"BEL":{"CN":"贝伦,巴西","EN":"Belém, Brazil"},"BEY":{"CN":"贝鲁特,黎巴嫩","EN":"Beirut, Lebanon"},"BGW":{"CN":"巴格达,伊拉克","EN":"Baghdad, Iraq"},"BKK":{"CN":"曼谷,泰国","EN":"Bangkok, Thailand"},"BLR":{"CN":"班加罗尔,印度","EN":"Bangalore, India"},"BNA":{"CN":"纳什维尔,田纳西州,美国","EN":"Nashville, TN, United States"},"BNE":{"CN":"布里斯班,昆士兰州,澳大利亚","EN":"Brisbane, QLD, Australia"},"BNU":{"CN":"布卢梅瑙,巴西","EN":"Blumenau, Brazil"},"BOG":{"CN":"波哥大,哥伦比亚","EN":"Bogotá, Colombia"},"BOM":{"CN":"孟买,印度","EN":"Mumbai, India"},"BOS":{"CN":"波士顿,马萨诸塞州,美国","EN":"Boston, MA, United States"},"BRU":{"CN":"布鲁塞尔,比利时","EN":"Brussels, Belgium"},"BSB":{"CN":"巴西利亚,巴西","EN":"Brasilia, Brazil"},"BUD":{"CN":"布达佩斯,匈牙利","EN":"Budapest, Hungary"},"BUF":{"CN":"水牛,纽约州,美国","EN":"Buffalo, NY, United States"},"BWN":{"CN":"斯里巴加湾市,文莱","EN":"Bandar Seri Begawan, Brunei"},"CAN":{"CN":"广州,中国","EN":"Guangzhou, China"},"CBR":{"CN":"堪培拉,澳大利亚首都领地,澳大利亚","EN":"Canberra, ACT, Australia"},"CCU":{"CN":"加尔各答,印度","EN":"Kolkata, India"},"CDG":{"CN":"巴黎,法国","EN":"Paris, France"},"CEB":{"CN":"宿雾,菲律宾","EN":"Cebu, Philippines"},"CFC":{"CN":"卡萨多,巴西","EN":"Caçador, Brazil"},"CGK":{"CN":"雅加达,印度尼西亚","EN":"Jakarta, Indonesia"},"CGO":{"CN":"郑州,中国","EN":"Zhengzhou, China"},"CGP":{"CN":"吉大港,孟加拉国","EN":"Chittagong, Bangladesh"},"CKG":{"CN":"重庆,中国","EN":"Chongqing, China"},"CLT":{"CN":"夏洛特,北卡罗来纳州,美国","EN":"Charlotte, NC, United States"},"CMB":{"CN":"科伦坡,斯里兰卡","EN":"Colombo, Sri Lanka"},"CMH":{"CN":"哥伦布,俄亥俄州,美国","EN":"Columbus, OH, United States"},"CMN":{"CN":"卡萨布兰卡,摩洛哥","EN":"Casablanca, Morocco"},"CNF":{"CN":"贝洛奥里藏特,巴西","EN":"Belo Horizonte, Brazil"},"CPH":{"CN":"哥本哈根,丹麦","EN":"Copenhagen, Denmark"},"CPT":{"CN":"开普敦,南非","EN":"Cape Town, South Africa"},"CSX":{"CN":"株洲,中国","EN":"Zhuzhou, China"},"CTU":{"CN":"成都,中国","EN":"Chengdu, China"},"CUR":{"CN":"威廉斯塔德,库拉索","EN":"Willemstad, Curaçao"},"CWB":{"CN":"库里蒂巴,巴西","EN":"Curitiba, Brazil"},"DAC":{"CN":"达卡,孟加拉国","EN":"Dhaka, Bangladesh"},"DAR":{"CN":"达累斯萨拉姆,坦桑尼亚","EN":"Dar Es Salaam, Tanzania"},"DEL":{"CN":"新德里,印度","EN":"New Delhi, India"},"DEN":{"CN":"丹佛,科罗拉多州,美国","EN":"Denver, CO, United States"},"DFW":{"CN":"达拉斯,德克萨斯州,美国","EN":"Dallas, TX, United States"},"DKR":{"CN":"达喀尔,塞内加尔","EN":"Dakar, Senegal"},"DME":{"CN":"莫斯科,俄罗斯","EN":"Moscow, Russia"},"DMM":{"CN":"达曼,沙特阿拉伯","EN":"Dammam, Saudi Arabia"},"DOH":{"CN":"多哈,卡塔尔","EN":"Doha, Qatar"},"DTW":{"CN":"底特律,密歇根州,美国","EN":"Detroit, MI, United States"},"DUB":{"CN":"都柏林,爱尔兰","EN":"Dublin, Ireland"},"DUR":{"CN":"德班,南非","EN":"Durban, South Africa"},"DUS":{"CN":"杜塞尔多夫,德国","EN":"Düsseldorf, Germany"},"DXB":{"CN":"迪拜,阿拉伯联合酋长国","EN":"Dubai, United Arab Emirates"},"EDI":{"CN":"爱丁堡,英国","EN":"Edinburgh, United Kingdom"},"EVN":{"CN":"埃里温,亚美尼亚","EN":"Yerevan, Armenia"},"EWR":{"CN":"纽瓦克,新泽西州,美国","EN":"Newark, NJ, United States"},"EZE":{"CN":"布宜诺斯艾利斯,阿根廷","EN":"Buenos Aires, Argentina"},"FCO":{"CN":"罗马,意大利","EN":"Rome, Italy"},"FLN":{"CN":"弗洛里亚诺波利斯,巴西","EN":"Florianopolis, Brazil"},"FOR":{"CN":"福塔莱萨,巴西","EN":"Fortaleza, Brazil"},"FRA":{"CN":"法兰克福,德国","EN":"Frankfurt, Germany"},"GIG":{"CN":"里约热内卢,巴西","EN":"Rio de Janeiro, Brazil"},"GND":{"CN":"圣乔治,格林纳达","EN":"St. George's, Grenada"},"GOT":{"CN":"哥德堡,瑞典","EN":"Gothenburg, Sweden"},"GRU":{"CN":"圣保罗,巴西","EN":"São Paulo, Brazil"},"GUA":{"CN":"危地马拉城,危地马拉","EN":"Guatemala City, Guatemala"},"GVA":{"CN":"日内瓦,瑞士","EN":"Geneva, Switzerland"},"GYD":{"CN":"巴库,阿塞拜疆","EN":"Baku, Azerbaijan"},"GYE":{"CN":"瓜亚基尔,厄瓜多尔","EN":"Guayaquil, Ecuador"},"HAM":{"CN":"汉堡,德国","EN":"Hamburg, Germany"},"HAN":{"CN":"河内,越南","EN":"Hanoi, Vietnam"},"HEL":{"CN":"赫尔辛基,芬兰","EN":"Helsinki, Finland"},"HKG":{"CN":"香港","EN":"Hong Kong"},"HNL":{"CN":"火奴鲁鲁,夏威夷州,美国","EN":"Honolulu, HI, United States"},"HRE":{"CN":"哈拉雷,津巴布韦","EN":"Harare, Zimbabwe"},"HYD":{"CN":"海得拉巴,印度","EN":"Hyderabad, India"},"IAD":{"CN":"阿什本,弗吉尼亚州,美国","EN":"Ashburn, VA, United States"},"IAH":{"CN":"休斯顿,德克萨斯州,美国","EN":"Houston, TX, United States"},"ICN":{"CN":"汉城,韩国","EN":"Seoul, South Korea"},"IND":{"CN":"印第安纳波利斯,印第安纳州,美国","EN":"Indianapolis, IN, United States"},"ISB":{"CN":"伊斯兰堡,巴基斯坦","EN":"Islamabad, Pakistan"},"IST":{"CN":"伊斯坦布尔,火鸡","EN":"Istanbul, Turkey"},"ITJ":{"CN":"伊塔雅伊,巴西","EN":"Itajaí, Brazil"},"JAX":{"CN":"杰克逊维尔,佛罗里达州,美国","EN":"Jacksonville, FL, United States"},"JHB":{"CN":"柔佛州新山,马来西亚","EN":"Johor Bahru, Malaysia"},"JIB":{"CN":"吉布提市,吉布提","EN":"Djibouti City, Djibouti"},"JNB":{"CN":"约翰内斯堡,南非","EN":"Johannesburg, South Africa"},"JSR":{"CN":"贾肖尔,孟加拉国","EN":"Jashore, Bangladesh"},"KBP":{"CN":"基辅,乌克兰","EN":"Kyiv, Ukraine"},"KEF":{"CN":"雷克雅未克,冰岛","EN":"Reykjavík, Iceland"},"KGL":{"CN":"基加利,卢旺达","EN":"Kigali, Rwanda"},"KHI":{"CN":"卡拉奇,巴基斯坦","EN":"Karachi, Pakistan"},"KIV":{"CN":"基希讷乌,摩尔多瓦","EN":"Chișinău, Moldova"},"KIX":{"CN":"大阪,日本","EN":"Osaka, Japan"},"KJA":{"CN":"克拉斯诺亚尔斯克,俄罗斯","EN":"Krasnoyarsk, Russia"},"KTM":{"CN":"加德满都,尼泊尔","EN":"Kathmandu, Nepal"},"KUL":{"CN":"吉隆坡,马来西亚","EN":"Kuala Lumpur, Malaysia"},"KWI":{"CN":"科威特城,科威特","EN":"Kuwait City, Kuwait"},"LAD":{"CN":"罗安达,安哥拉","EN":"Luanda, Angola"},"LAS":{"CN":"拉斯维加斯,内华达州,美国","EN":"Las Vegas, NV, United States"},"LAX":{"CN":"洛杉矶,加利福尼亚州,美国","EN":"Los Angeles, CA, United States"},"LCA":{"CN":"尼科西亚,塞浦路斯","EN":"Nicosia, Cyprus"},"LED":{"CN":"圣彼得堡,俄罗斯","EN":"Saint Petersburg, Russia"},"LHE":{"CN":"拉合尔,巴基斯坦","EN":"Lahore, Pakistan"},"LHR":{"CN":"伦敦,英国","EN":"London, United Kingdom"},"LIM":{"CN":"利马,秘鲁","EN":"Lima, Peru"},"LIS":{"CN":"里斯本,葡萄牙","EN":"Lisbon, Portugal"},"LOS":{"CN":"拉各斯,尼日利亚","EN":"Lagos, Nigeria"},"LUX":{"CN":"卢森堡城,卢森堡","EN":"Luxembourg City, Luxembourg"},"MAA":{"CN":"钦奈,印度","EN":"Chennai, India"},"MAD":{"CN":"马德里,西班牙","EN":"Madrid, Spain"},"MAN":{"CN":"曼彻斯特,英国","EN":"Manchester, United Kingdom"},"MBA":{"CN":"蒙巴萨,肯尼亚","EN":"Mombasa, Kenya"},"MCI":{"CN":"堪萨斯城,密苏里州,美国","EN":"Kansas City, MO, United States"},"MCT":{"CN":"马斯喀特,阿曼","EN":"Muscat, Oman"},"MDE":{"CN":"麦德林,哥伦比亚","EN":"Medellín, Colombia"},"MEL":{"CN":"墨尔本,维多利亚州,澳大利亚","EN":"Melbourne, VIC, Australia"},"MEM":{"CN":"孟菲斯,田纳西州,美国","EN":"Memphis, TN, United States"},"MEX":{"CN":"墨西哥城,墨西哥","EN":"Mexico City, Mexico"},"MFE":{"CN":"麦卡伦,德克萨斯州,美国","EN":"McAllen, TX, United States"},"MFM":{"CN":"澳门","EN":"Macau"},"MGM":{"CN":"蒙哥马利,亚拉巴马州,美国","EN":"Montgomery, AL, United States"},"MIA":{"CN":"迈阿密,佛罗里达州,美国","EN":"Miami, FL, United States"},"MLE":{"CN":"男性,马尔代夫","EN":"Malé, Maldives"},"MNL":{"CN":"马尼拉,菲律宾","EN":"Manila, Philippines"},"MPM":{"CN":"马普托,莫桑比克","EN":"Maputo, Mozambique"},"MRS":{"CN":"马赛,法国","EN":"Marseille, France"},"MRU":{"CN":"路易港,毛里求斯","EN":"Port Louis, Mauritius"},"MSP":{"CN":"明尼阿波利斯,明尼苏达州,美国","EN":"Minneapolis, MN, United States"},"MUC":{"CN":"慕尼黑,德国","EN":"Munich, Germany"},"MXP":{"CN":"米兰,意大利","EN":"Milan, Italy"},"NAG":{"CN":"那格浦尔,印度","EN":"Nagpur, India"},"NBG":{"CN":"宁波,中国","EN":"Ningbo, China"},"NBO":{"CN":"内罗毕,肯尼亚","EN":"Nairobi, Kenya"},"NOU":{"CN":"努美阿,新喀里多尼亚","EN":"Noumea, New Caledonia"},"NRT":{"CN":"东京,日本","EN":"Tokyo, Japan"},"OMA":{"CN":"奥马哈,内布拉斯加州,美国","EN":"Omaha, NE, United States"},"ORD":{"CN":"芝加哥,伊利诺伊州,美国","EN":"Chicago, IL, United States"},"ORF":{"CN":"诺福克,弗吉尼亚州,美国","EN":"Norfolk, VA, United States"},"ORK":{"CN":"软木,爱尔兰","EN":"Cork, Ireland"},"OSL":{"CN":"奥斯陆,挪威","EN":"Oslo, Norway"},"OTP":{"CN":"布加勒斯特,罗马尼亚","EN":"Bucharest, Romania"},"PAP":{"CN":"太子港,海地","EN":"Port-Au-Prince, Haiti"},"PBH":{"CN":"廷布,不丹","EN":"Thimphu, Bhutan"},"PBM":{"CN":"帕拉马里博,苏里南","EN":"Paramaribo, Suriname"},"PDX":{"CN":"波特兰,俄勒冈州,美国","EN":"Portland, OR, United States"},"PER":{"CN":"珀斯,西澳州,澳大利亚","EN":"Perth, WA, Australia"},"PHL":{"CN":"费城,美国","EN":"Philadelphia, United States"},"PHX":{"CN":"凤凰,亚利桑那州,美国","EN":"Phoenix, AZ, United States"},"PIT":{"CN":"匹兹堡,宾夕法尼亚州,美国","EN":"Pittsburgh, PA, United States"},"PMO":{"CN":"巴勒莫,意大利","EN":"Palermo, Italy"},"PNH":{"CN":"金边,柬埔寨","EN":"Phnom Penh, Cambodia"},"POA":{"CN":"阿雷格里港,巴西","EN":"Porto Alegre, Brazil"},"PRG":{"CN":"布拉格,捷克共和国","EN":"Prague, Czech Republic"},"PTY":{"CN":"巴拿马城,巴拿马","EN":"Panama City, Panama"},"QRO":{"CN":"克雷塔罗,墨西哥,墨西哥","EN":"Queretaro, MX, Mexico"},"QWJ":{"CN":"美洲,巴西","EN":"Americana, Brazil"},"RAO":{"CN":"里贝朗普雷图,巴西","EN":"Ribeirao Preto, Brazil"},"RGN":{"CN":"仰光,缅甸","EN":"Yangon, Myanmar"},"RIC":{"CN":"里士满,弗吉尼亚州,美国","EN":"Richmond, VA, United States"},"RIX":{"CN":"里加,拉脱维亚","EN":"Riga, Latvia"},"ROB":{"CN":"蒙罗维亚,利比里亚","EN":"Monrovia, Liberia"},"RUH":{"CN":"利雅得,沙特阿拉伯","EN":"Riyadh, Saudi Arabia"},"RUN":{"CN":"团圆,法国","EN":"Réunion, France"},"SAN":{"CN":"圣地亚哥,加利福尼亚州,美国","EN":"San Diego, CA, United States"},"SCL":{"CN":"圣地亚哥,智利","EN":"Santiago, Chile"},"SEA":{"CN":"西雅图,华盛顿州,美国","EN":"Seattle, WA, United States"},"SGN":{"CN":"胡志明市,越南","EN":"Ho Chi Minh City, Vietnam"},"SHA":{"CN":"上海,中国","EN":"Shanghai, China"},"SIN":{"CN":"新加坡,新加坡","EN":"Singapore, Singapore"},"SJC":{"CN":"圣荷西,加利福尼亚州,美国","EN":"San Jose, CA, United States"},"SJO":{"CN":"圣荷西,哥斯达黎加","EN":"San José, Costa Rica"},"SJP":{"CN":"圣何塞-杜里约普雷图,巴西","EN":"São José do Rio Preto, Brazil"},"SKG":{"CN":"塞萨洛尼基,希腊","EN":"Thessaloniki, Greece"},"SLC":{"CN":"盐湖城,犹他州,美国","EN":"Salt Lake City, UT, United States"},"SMF":{"CN":"萨克拉门托,加利福尼亚州,美国","EN":"Sacramento, CA, United States"},"SOD":{"CN":"索罗卡巴,巴西","EN":"Sorocaba, Brazil"},"SOF":{"CN":"苏菲亚,保加利亚","EN":"Sofia, Bulgaria"},"SSA":{"CN":"萨尔瓦多,巴西","EN":"Salvador, Brazil"},"STL":{"CN":"圣路易斯,密苏里州,美国","EN":"St. Louis, MO, United States"},"SVX":{"CN":"叶卡捷琳堡,俄罗斯","EN":"Yekaterinburg, Russia"},"SYD":{"CN":"悉尼,新南威尔士州,澳大利亚","EN":"Sydney, NSW, Australia"},"SZV":{"CN":"苏州,中国","EN":"Suzhou, China"},"TBS":{"CN":"第比利斯,乔治亚州","EN":"Tbilisi, Georgia"},"TGU":{"CN":"特古西加尔巴,洪都拉斯","EN":"Tegucigalpa, Honduras"},"TLH":{"CN":"塔拉哈西,佛罗里达州,美国","EN":"Tallahassee, FL, United States"},"TLL":{"CN":"塔林,爱沙尼亚","EN":"Tallinn, Estonia"},"TLV":{"CN":"特拉维夫,以色列","EN":"Tel Aviv, Israel"},"TNA":{"CN":"济南,中国","EN":"Jinan, China"},"TNR":{"CN":"塔那那利佛,马达加斯加","EN":"Antananarivo, Madagascar"},"TPA":{"CN":"坦帕,佛罗里达州,美国","EN":"Tampa, FL, United States"},"TPE":{"CN":"台北","EN":"Taipei"},"TSN":{"CN":"天津,中国","EN":"Tianjin, China"},"TUN":{"CN":"突尼斯,突尼斯","EN":"Tunis, Tunisia"},"TXL":{"CN":"柏林,德国","EN":"Berlin, Germany"},"UIO":{"CN":"基多,厄瓜多尔","EN":"Quito, Ecuador"},"ULN":{"CN":"乌兰巴托,蒙古","EN":"Ulaanbaatar, Mongolia"},"URT":{"CN":"素叻他尼,泰国","EN":"Surat Thani, Thailand"},"VCP":{"CN":"坎皮纳斯,巴西","EN":"Campinas, Brazil"},"VIE":{"CN":"维也纳,奥地利","EN":"Vienna, Austria"},"VNO":{"CN":"维尔纽斯,立陶宛","EN":"Vilnius, Lithuania"},"VTE":{"CN":"万象,老挝","EN":"Vientiane, Laos"},"WAW":{"CN":"华沙,波兰","EN":"Warsaw, Poland"},"WUH":{"CN":"武汉,中国","EN":"Wuhan, China"},"WUX":{"CN":"无锡,中国","EN":"Wuxi, China"},"XIY":{"CN":"西安,中国","EN":"Xi'an, China"},"YUL":{"CN":"蒙特利尔,魁北克,加拿大","EN":"Montréal, QC, Canada"},"YVR":{"CN":"温哥华,不列颠哥伦比亚,加拿大","EN":"Vancouver, BC, Canada"},"YWG":{"CN":"温尼伯,曼尼托巴,加拿大","EN":"Winnipeg, MB, Canada"},"YXE":{"CN":"萨斯卡通,萨斯喀彻温,加拿大","EN":"Saskatoon, SK, Canada"},"YYC":{"CN":"卡尔加里,艾伯塔,加拿大","EN":"Calgary, AB, Canada"},"YYZ":{"CN":"多伦多,安大略,加拿大","EN":"Toronto, ON, Canada"},"ZAG":{"CN":"萨格勒布,克罗地亚","EN":"Zagreb, Croatia"},"ZDM":{"CN":"拉马拉","EN":"Ramallah"},"ZRH":{"CN":"苏黎世,瑞士","EN":"Zürich, Switzerland"}}`
)

func GetColoInfoEN(colo string) string {
	colomap := make(map[string]map[string]string)
	err := json.Unmarshal([]byte(coloinfo), &colomap)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if val, ok := colomap[colo]; ok {
		return val["EN"]
	}
	return ""
}

func main() {
	fmt.Println("speed_cloudflare client version:1.0.0 https://github.com/bianzhifu/speedcf")
	flag.IntVar(&ipv, "ipv", 0, "网络:0、默认 4、IPV4 6、IPV6")
	flag.Parse()
	if ipv == 0 {
		SetTransportProtocol("tcp")
	} else if ipv == 4 {
		SetTransportProtocol("tcp4")
	} else if ipv == 6 {
		SetTransportProtocol("tcp6")
	} else {
		fmt.Println("网络选择错误")
		os.Exit(0)
	}

	meta := GetMetadata()
	fmt.Println(fmt.Sprintf("     Server location : %s (%s)", meta.CF_Mete_Colo, meta.CF_Mete_ColoInfo))
	fmt.Println(fmt.Sprintf("     Client location : %s,%s", meta.CF_Mete_Country, meta.CF_Mete_City))
	fmt.Println(fmt.Sprintf("           Client IP : %s (AS%s)", meta.CF_Mete_IP, meta.CF_Mete_ASN))
	ping, err := Ping()
	if err != nil {
		fmt.Println(fmt.Sprintf("             Latency : %s", err.Error()))
	} else {
		fmt.Println(fmt.Sprintf("             Latency : %dms", ping.Milliseconds()))

	}
	m := 5
	codes := []string{"100kB", "  1MB", " 10MB", " 25MB", "100MB"}
	sizes := []int64{100 * 1024, 1024 * 1024, 10 * 1024 * 1024, 25 * 1024 * 1024, 100 * 1024 * 1024}
	for i, size := range sizes {
		downresult := SpeedFunc("down", size, m)
		for i := len(downresult); i < 7; i++ {
			downresult = " " + downresult
		}
		upresult := SpeedFunc("up", size, m)
		for i := len(upresult); i < 7; i++ {
			upresult = " " + upresult
		}
		fmt.Println("        " + codes[i] + "  Speed : ↓ " + downresult + "Mbps   ↑ " + upresult + "Mbps")
	}
}

func SetTransportProtocol(protocol string) {
	http.DefaultTransport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: func(ctx context.Context, _, addr string) (net.Conn, error) {
			return (&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext(ctx, protocol, addr)
		},
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true}, //忽略证书
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}

type Metadata struct {
	CF_Mete_IP       string
	CF_Mete_ASN      string
	CF_Mete_City     string
	CF_Mete_Country  string
	CF_Mete_Colo     string
	CF_Mete_ColoInfo string
}

func GetMetadata() *Metadata {
	resurlt, err := Downlink(0)
	if err != nil {
		time.Sleep(time.Second)
		return GetMetadata()
	}
	return &Metadata{
		CF_Mete_IP:       resurlt.HTTPRespHeader.Get("cf-meta-ip"),
		CF_Mete_ASN:      resurlt.HTTPRespHeader.Get("cf-meta-asn"),
		CF_Mete_City:     resurlt.HTTPRespHeader.Get("cf-meta-city"),
		CF_Mete_Country:  resurlt.HTTPRespHeader.Get("cf-meta-country"),
		CF_Mete_Colo:     resurlt.HTTPRespHeader.Get("cf-meta-colo"),
		CF_Mete_ColoInfo: GetColoInfoEN(resurlt.HTTPRespHeader.Get("cf-meta-colo")),
	}
}

func Ping() (time.Duration, error) {
	durationSum := time.Duration(0)
	rttCount := 0
	for start := time.Now(); time.Since(start) < 2*time.Second; {
		measurement, err := Downlink(0)
		if err != nil {
			return time.Duration(0), err
		}
		cfReqDurMatch := regexp.MustCompile(`cfRequestDuration;dur=([\d.]+)`).FindStringSubmatch(measurement.HTTPRespHeader.Get("Server-Timing"))
		if len(cfReqDurMatch) > 0 {
			cfReqDur, err := time.ParseDuration(fmt.Sprintf("%sms", cfReqDurMatch[1]))
			if err == nil {
				durationSum += measurement.Duration - cfReqDur
				rttCount += 1
			}
		}
	}
	if rttCount == 0 {
		return time.Duration(0), errors.New("Rtt fail")
	}
	return durationSum / time.Duration(rttCount), nil
}

func SpeedFunc(mode string, size int64, m int) string {
	allBytes := int64(0)
	allDuration := int64(0)
	for j := 0; j < m; j++ {
		var result *SpeedResult
		switch mode {
		case "down":
			result, _ = Downlink(size)
		case "up":
			result, _ = Uplink(size)
		}
		atomic.AddInt64(&allBytes, result.Size)
		atomic.AddInt64(&allDuration, int64(result.Duration))
	}
	bw := float64(allBytes*8) / time.Duration(allDuration).Seconds()
	return fmt.Sprintf("%.2f", float64(bw)/float64(1024*1024))
}

type SpeedResult struct {
	Size           int64
	Duration       time.Duration
	HTTPRespHeader http.Header
}

func Downlink(size int64) (*SpeedResult, error) {
	start := time.Now()
	resp, err := http.Get("https://speed.cloudflare.com/__down?bytes=" + strconv.FormatInt(size, 10))
	if err != nil {
		return nil, err
	}
	downloadedSize, err := io.Copy(ioutil.Discard, resp.Body)
	if err != nil {
		return nil, err
	}
	end := time.Now()
	return &SpeedResult{
		Size:           downloadedSize,
		Duration:       end.Sub(start),
		HTTPRespHeader: resp.Header,
	}, nil
}

func Uplink(size int64) (*SpeedResult, error) {
	postBodyReader := InitReadMock(size)
	start := time.Now()
	resp, err := http.Post("https://speed.cloudflare.com/__up", "application/octet-stream", postBodyReader)
	if err != nil {
		return nil, err
	}

	end := time.Now()

	_, err = io.Copy(ioutil.Discard, resp.Body)
	if err != nil {
		return nil, err
	}
	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return &SpeedResult{
		Size:           size,
		Duration:       end.Sub(start),
		HTTPRespHeader: resp.Header,
	}, nil
}

type ReadMock struct {
	cidx int64
	cEOF int64
}

func (r *ReadMock) Read(p []byte) (int, error) {
	var err error = nil
	size := len(p)
	size64 := int64(size)
	r.cidx += size64
	if r.cidx >= r.cEOF {
		size = int(size64 - (r.cidx - r.cEOF))
		err = io.EOF
		r.cidx = r.cEOF
	}
	return size, err
}

func InitReadMock(size int64) *ReadMock {
	r := &ReadMock{}
	r.cidx = 0
	r.cEOF = size
	return r
}
