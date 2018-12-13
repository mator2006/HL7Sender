package main

import (
	"bufio"
	"fmt"
	"gopkg.in/gcfg.v1"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {

	var Host,
		Port,
		hl7mesg,
		hl7mesgsub,
		PatientNAME,
		PatientBD,
		PatientSEX,
		PatientID,
		OrderControl,
		AccessionNO,
		RPNO, SPNO,
		MODALITY,
		SPSDESC,
		StationNAME,
		StationAET,
		configurefile,
		StudyInstanceUID	string

	configurefile = "config.ini"
	iniexist, err := os.Stat(configurefile)

	if iniexist != nil {
		fmt.Println("Found Configure file.\nConfigure file is [" + configurefile + "].")
		config := struct {
			Netconn struct {
				Host string
				Port string
			}
			HL7order struct {
				PatientID    		string
				PatientNAME  		string
				PatientBD    		string
				PatientSEX   		string
				OrderControl 		string
				AccessionNO  		string
				RPNO         		string
				SPNO         		string
				MODALITY     		string
				SPSDESC      		string
				StationNAME  		string
				StationAET   		string
				StudyInstanceUID	string
			}
		}{}

		err = gcfg.ReadFileInto(&config, configurefile)
		if err != nil {
			fmt.Println("Failed to parse[fy:fen xi] configure file:")
			fmt.Println(err)
			fmt.Println("You can DO THIS:\n1.Tyr modify the configure file:[" + configurefile + "]\n2.Or remove this configure files,Will use DEFAULT configure")
			return
		} else {
			Host = strings.TrimSpace(config.Netconn.Host)
			Port = strings.TrimSpace(config.Netconn.Port)

			PatientID = strings.TrimSpace(config.HL7order.PatientID)
			PatientNAME = strings.Replace(strings.TrimSpace(config.HL7order.PatientNAME), " ", "^", -1)
			PatientBD = strings.TrimSpace(config.HL7order.PatientBD)
			PatientSEX = strings.TrimSpace(config.HL7order.PatientSEX)
			OrderControl = strings.TrimSpace(config.HL7order.OrderControl)
			AccessionNO = strings.TrimSpace(config.HL7order.AccessionNO)
			RPNO = strings.TrimSpace(config.HL7order.RPNO)
			SPNO = strings.TrimSpace(config.HL7order.SPNO)
			MODALITY = strings.TrimSpace(config.HL7order.MODALITY)
			SPSDESC = strings.TrimSpace(config.HL7order.SPSDESC)
			StationNAME = strings.TrimSpace(config.HL7order.StationNAME)
			StationAET = strings.TrimSpace(config.HL7order.StationAET)
			StudyInstanceUID = strings.TrimSpace(config.HL7order.StudyInstanceUID)
			
		}
	} else {
		fmt.Println("Not found configure file !\nUse DEFAULT config!")

		fmt.Println("Please enter HL7 server IP address:")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		Host = input.Text()

		if Host == "" || net.ParseIP(Host) == nil {
			fmt.Println("Invalid IP address.\nEnd.")
			return
		}

		fmt.Println("Please enter HL7 server PORT[Default 2575]:")
		input = bufio.NewScanner(os.Stdin)
		input.Scan()
		Port = input.Text()

		Portd := false
		if Port != "" {
			Porttoi, err := strconv.Atoi(Port)
			if err == nil {
				if Porttoi > 0 && Porttoi < 65536 {
					Port = strconv.Itoa(Porttoi)
					Portd = true
				}
			}
		}
		if !Portd {
			fmt.Println("Invalid Port number\nUse Default Port [2575]")
			Port = "2575"
		}
		_ = Portd

		rand.Seed(time.Now().Unix())
		DATESN := time.Now().Format("20060102")
		QUEENSN := time.Now().Format("150405") + "0" + strconv.Itoa(rand.Intn(10))

		PatientID = DATESN + QUEENSN
		PatientNAME = "Liu^Bei"
		PatientBD = "19840420"
		PatientSEX = "M"
		OrderControl = "NW"
		AccessionNO = DATESN + QUEENSN
		RPNO = "RP" + AccessionNO[2:len(DATESN + QUEENSN)]
		SPNO = "SP" + AccessionNO[2:len(DATESN + QUEENSN)]
		MODALITY = "MR"
		SPSDESC = "Tou Lu"
		StationNAME = "NO3RF"
		StationAET = "AW44"
		StudyInstanceUID = "2.16.840.1.113929.1.2.6493."+DATESN+"."+QUEENSN+".762142" //https://blog.csdn.net/dragonlet/article/details/37052997
	}

	_ = iniexist
	//_=jsonexist
	//_=yamlexist
	hl7mesgsub = //此处修改需谨慎
		`MSH|^~\&|||||||ORM^O01||P|2.3|||||CN` + "\r" +
			`PID|||` + PatientID + `^^^HIS||` + PatientNAME + `||` + PatientBD + `|` + PatientSEX + `|||||||||||||||||||||||||||||||` + "\r" +
			`PV1||E|||||||||||||||||V103-1^^^ADT1||||||||||||||||||||||||||||||||V|` + "\r" +
			`ORC|` + OrderControl + `||||SC||1^once^^^^S||||||||||||` + "\r" +
			`OBR|1|||^^^^` + SPSDESC + `||||||||||||||` + AccessionNO + `|` + RPNO + `|` + SPNO + `||||` + MODALITY + `|||1^once^^^^S|||WALK||||||||||||||P1^` + SPSDESC + `^ERL_MESA, ZDS|1.113654.3.13.1025^100^Application^DICOM` + "\r" +
			`ZDS|`+StudyInstanceUID+`^MEPACS^Application^DICOM|` + StationNAME + `|` + StationAET + "\r"
	hl7mesg = "\v" + hl7mesgsub + "\x1C" + "\r" //谨慎同上

	//net.dial 拨号 获取tcp连接
	constr := Host + ":" + Port
	conn, err := net.Dial("tcp", constr)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("获取[" + constr + "]的tcp连接成功.")

	//发送HL7数据
	sendbytecount, err := fmt.Fprintf(conn, hl7mesg)
	if err != nil {
		fmt.Println("发送错误！", err)
		return
	}
	fmt.Println("共发送", sendbytecount, "字节.\n发送hl7数据完成.")

	//接收返回数据
	Receivemesg, _ := bufio.NewReader(conn).ReadString('\x1C') //读取返回数据，遇到指定字符结束

	if !strings.Contains(string(Receivemesg), `AR`) { //不知道为什么，直接写包含AA就是不行
		fmt.Println("传输成功.")
	} else {
		fmt.Println("失败，未知错误!")
		fmt.Println(string(Receivemesg))
	}
	//关闭连接
	defer conn.Close()
	fmt.Println("连接已关闭.")
}

//OBR|1|||| to OBR|1|||^^^^`+SPSDESC+`| mator 2018.12.13
// test git commit 5
