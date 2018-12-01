package main

import (
	"fmt"
	"net"
	"bufio"
	"strings"
	"time"
	"math/rand"
	"strconv"
	"os"
	"gopkg.in/gcfg.v1"
)

var (
	hl7mesgsub string 
	Host string 
	Port string
)

func main() {

	configurefile := ""
	iniexist,err := os.Stat ("config.ini")
	jsonexist,err := os.Stat ("config.json")
	yamlexist,err := os.Stat ("config.yaml")

	if (iniexist != nil || jsonexist != nil|| yamlexist != nil){
		if iniexist != nil {
			configurefile = "config.ini"
		}
		if jsonexist != nil {
			configurefile = "config.json"
		}
		if yamlexist != nil {
			configurefile = "config.yaml"
		}
		fmt.Println("Found Configure file.\nConfigure file is ["+configurefile+"].")
		config := struct {
			Netconn struct{
				Host string
				Port string
			}
			HL7order struct{
				PatientIDID string
				PatientIDNAME string
				PatientIDBD string
				PatientIDSEX string
				OrderControl string
				AccessionNO string
				RPNO string
				SPNO string
				MODALITY string
				SPSDESC string
				StationNAME string
				StationAET string
			}
		}{}

		err = gcfg.ReadFileInto(&config,configurefile)
			if err != nil {
				fmt.Println("Failed to parse[fy:fen xi] configure file:")
				fmt.Println(err)
				fmt.Println("You can DO THIS:\n1.Tyr modify the configure file:["+configurefile+"]\n2.Or remove this configure files,Will use DEFAULT configure")
				return
			}else{
				Host = config.Netconn.Host
				Port = config.Netconn.Port

				PatientIDID := config.HL7order.PatientIDID
				PatientIDNAME := config.HL7order.PatientIDNAME
				PatientIDBD := config.HL7order.PatientIDBD
				PatientIDSEX := config.HL7order.PatientIDSEX
				OrderControl := config.HL7order.OrderControl
				AccessionNO := config.HL7order.AccessionNO
				RPNO := config.HL7order.RPNO
				SPNO := config.HL7order.SPNO
				MODALITY := config.HL7order.MODALITY
				SPSDESC := config.HL7order.SPSDESC
				StationNAME := config.HL7order.StationNAME
				StationAET := config.HL7order.StationAET

				hl7mesgsub =   //此处修改需谨慎
				`MSH|^~\&|||||||ORM^O01||P|2.3|||||CN`+ "\r" +
				`PID|||`+PatientIDID+`^^^HL7||`+PatientIDNAME+`^3||`+PatientIDBD+`|`+PatientIDSEX+`|||||||||||||||||||||||||||||||` + "\r" +
				`PV1||E|||||||||||||||||V103-1^^^ADT1||||||||||||||||||||||||||||||||V|`+ "\r" +
				`ORC|`+OrderControl+`||||SC||1^once^^^^S||||||||||||`+ "\r" +
				`OBR|1|||||||||||||||||`+AccessionNO+`|`+RPNO+`|`+SPNO+`||||`+MODALITY+`|||1^once^^^^S|||WALK||||||||||||||P1^`+SPSDESC+`^ERL_MESA, ZDS|1.113654.3.13.1025^100^Application^DICOM` + "\r" + 
				`ZDS|2.16.840.1.113929.1.2.6493.20070508.21948.762142^DCM4CHEE^Application^DICOM|`+StationNAME+`|`+StationAET+ "\r"
			}
		}else{
			fmt.Println("Not found configure file !\nUse DEFAULT config!!!")

			Host = "192.168.1.203"
			Port = "2575"

			rand.Seed(time.Now().Unix())
			DATESN := time.Now().Format("20060102")
			QUEENSN := time.Now().Format("150405")+"00"+strconv.Itoa(rand.Intn(10))

			PatientIDID := DATESN+QUEENSN
			PatientIDNAME := "Liu^Bei"
			PatientIDBD := "19840420"
			PatientIDSEX := "M"
			OrderControl := "NW"
			AccessionNO := DATESN+QUEENSN
			RPNO := "RP"+DATESN+QUEENSN
			SPNO := "SP"+DATESN+QUEENSN
			MODALITY := "MR"
			SPSDESC := "Tou Lu"
			StationNAME := "NO3RF"
			StationAET := "AW44"

			hl7mesgsub =   //此处修改需谨慎
			`MSH|^~\&|||||||ORM^O01||P|2.3|||||CN`+ "\r" +
			`PID|||`+PatientIDID+`^^^HL7||`+PatientIDNAME+`^3||`+PatientIDBD+`|`+PatientIDSEX+`|||||||||||||||||||||||||||||||` + "\r" +
			`PV1||E|||||||||||||||||V103-1^^^ADT1||||||||||||||||||||||||||||||||V|`+ "\r" +
			`ORC|`+OrderControl+`||||SC||1^once^^^^S||||||||||||`+ "\r" +
			`OBR|1|||||||||||||||||`+AccessionNO+`|`+RPNO+`|`+SPNO+`||||`+MODALITY+`|||1^once^^^^S|||WALK||||||||||||||P1^`+SPSDESC+`^ERL_MESA, ZDS|1.113654.3.13.1025^100^Application^DICOM` + "\r" + 
			`ZDS|2.16.840.1.113929.1.2.6493.20070508.21948.762142^DCM4CHEE^Application^DICOM|`+StationNAME+`|`+StationAET+ "\r"

		}

	_=iniexist
	_=jsonexist
	_=yamlexist

	hl7mesg := 	"\v" + hl7mesgsub + "\x1C" + "\r"  //谨慎同上

	//net.dial 拨号 获取tcp连接
	constr := Host + ":" + Port
	conn, err := net.Dial("tcp", constr)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("获取[" + constr + "]的tcp连接成功...")


	//发送HL7数据
	fmt.Fprintf(conn, hl7mesg)
	fmt.Println("发送hl7数据完成...")


	//接收返回数据
	Receivemesg, _:= bufio.NewReader(conn).ReadString('\x1C') //读取返回数据，遇到指定字符结束
	
	if !strings.Contains(string(Receivemesg), `AR`)  {   //不知道为什么，直接写包含AA就是不行
		fmt.Println("传输成功")
		}else{
		fmt.Println("失败，未知错误")
		fmt.Println( string(Receivemesg) )
	}
	//关闭连接
	defer conn.Close()
	fmt.Println("连接已关闭")
}
