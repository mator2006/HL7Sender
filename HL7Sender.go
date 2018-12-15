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

type hl7 struct{
	Host				string
	Port				string
	hl7mesg				string
	hl7mesgsub			string
	PatientNAME			string
	PatientBD			string
	PatientSEX			string
	PatientID			string
	OrderControl		string
	AccessionNO			string
	RPNO, SPNO			string
	MODALITY			string
	SPSDESC				string
	StationNAME			string
	StationAET			string
	configurefile		string
	StudyInstanceUID	string
}

func main() {

	var hl7data	hl7

	configurefile := "config.ini"
	iniexist, err := os.Stat(configurefile)

	if iniexist != nil {
		useconfigfile(hl7data,configurefile)
	} else {
		usedefconfig(hl7data)
	}
	
	_=err
	_ = iniexist

	addhl7data(hl7data)
	senddata(hl7data)
}

func useconfigfile(hl7data hl7,configurefile string){
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

	err := gcfg.ReadFileInto(&config, configurefile)
	if err != nil {
		fmt.Println("Failed to parse[fy:fen xi] configure file:")
		fmt.Println(err)
		fmt.Println("You can DO THIS:\n1.Tyr modify the configure file:[" + configurefile + "]\n2.Or remove this configure files,Will use DEFAULT configure")
		return
		} else {
			hl7data.Host = strings.TrimSpace(config.Netconn.Host)
			hl7data.Port = strings.TrimSpace(config.Netconn.Port)

			hl7data.PatientID = strings.TrimSpace(config.HL7order.PatientID)
			hl7data.PatientNAME = strings.Replace(strings.TrimSpace(config.HL7order.PatientNAME), " ", "^", -1)
			hl7data.PatientBD = strings.TrimSpace(config.HL7order.PatientBD)
			hl7data.PatientSEX = strings.TrimSpace(config.HL7order.PatientSEX)
			hl7data.OrderControl = strings.TrimSpace(config.HL7order.OrderControl)
			hl7data.AccessionNO = strings.TrimSpace(config.HL7order.AccessionNO)
			hl7data.RPNO = strings.TrimSpace(config.HL7order.RPNO)
			hl7data.SPNO = strings.TrimSpace(config.HL7order.SPNO)
			hl7data.MODALITY = strings.TrimSpace(config.HL7order.MODALITY)
			hl7data.SPSDESC = strings.TrimSpace(config.HL7order.SPSDESC)
			hl7data.StationNAME = strings.TrimSpace(config.HL7order.StationNAME)
			hl7data.StationAET = strings.TrimSpace(config.HL7order.StationAET)
			hl7data.StudyInstanceUID = strings.TrimSpace(config.HL7order.StudyInstanceUID)
			}
	}

func usedefconfig(hl7data hl7){
	fmt.Println("Not found configure file !\nUse DEFAULT config!")

	fmt.Println("Please enter HL7 server IP address:")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	vHost := input.Text()

	if vHost == "" || net.ParseIP(vHost) == nil {
		fmt.Println("Invalid IP address.\nEnd.")
		return
	}else{
		hl7data.Host = vHost
	}

	fmt.Println("Please enter HL7 server PORT[Default 2575]:")
	input = bufio.NewScanner(os.Stdin)
	input.Scan()
	vPort := input.Text()

	Portd := false
	if vPort != "" {
		Porttoi, err := strconv.Atoi(vPort)
		if err == nil {
			if Porttoi > 0 && Porttoi < 65536 {
				hl7data.Port = strconv.Itoa(Porttoi)
				Portd = true
			}
		}
	}
	if !Portd {
		fmt.Println("Invalid Port number\nUse Default Port [2575]")
		hl7data.Port = "2575"
	}
	_ = Portd

	rand.Seed(time.Now().Unix())
	DATESN := time.Now().Format("20060102")
	QUEENSN := time.Now().Format("150405") + "0" + strconv.Itoa(rand.Intn(10))

	hl7data.PatientID = DATESN + QUEENSN
	hl7data.PatientNAME = "Liu^Bei"
	hl7data.PatientBD = "19840420"
	hl7data.PatientSEX = "M"
	hl7data.OrderControl = "NW"
	hl7data.AccessionNO = DATESN + QUEENSN
	hl7data.RPNO = "RP" + hl7data.AccessionNO[2:len(DATESN + QUEENSN)]
	hl7data.SPNO = "SP" + hl7data.AccessionNO[2:len(DATESN + QUEENSN)]
	hl7data.MODALITY = "MR"
	hl7data.SPSDESC = "Tou Lu"
	hl7data.StationNAME = "NO3RF"
	hl7data.StationAET = "AW44"
	hl7data.StudyInstanceUID = "2.16.840.1.113929.1.2.6493."+DATESN+"."+QUEENSN+".762142" //https://blog.csdn.net/dragonlet/article/details/37052997
}

func addhl7data(hl7data hl7){
	hl7data.hl7mesgsub = `MSH|^~\&|||||||ORM^O01||P|2.3|||||CN`+"\r"+`PID|||` + hl7data.PatientID + `^^^HIS||`+hl7data.PatientNAME + `||` +hl7data.PatientBD + `|` + hl7data.PatientSEX + `|||||||||||||||||||||||||||||||` + "\r" +
		`PV1||E|||||||||||||||||V103-1^^^ADT1||||||||||||||||||||||||||||||||V|` + "\r" +
		`ORC|` + hl7data.OrderControl + `||||SC||1^once^^^^S||||||||||||` + "\r" +
		`OBR|1|||^^^^` + hl7data.SPSDESC + `||||||||||||||` + hl7data.AccessionNO + `|` + hl7data.RPNO + `|` + hl7data.SPNO + `||||` + hl7data.MODALITY +
		 `|||1^once^^^^S|||WALK||||||||||||||P1^` + hl7data.SPSDESC + `^ERL_MESA, ZDS|1.113654.3.13.1025^100^Application^DICOM` + "\r" +
		`ZDS|`+hl7data.StudyInstanceUID+`^MEPACS^Application^DICOM|` + hl7data.StationNAME + `|` + hl7data.StationAET + "\r"
		hl7data.hl7mesg = "\v" + hl7data.hl7mesgsub + "\x1C" + "\r" //谨慎同上
	fmt.Println(hl7data.hl7mesgsub)
}

func senddata (hl7data hl7){
	//net.dial 拨号 获取tcp连接
	constr := hl7data.Host + ":" + hl7data.Port
	conn, err := net.Dial("tcp", constr)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("获取[" + constr + "]的tcp连接成功.")

	//发送HL7数据
	sendbytecount, err := fmt.Fprintf(conn, hl7data.hl7mesg)
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
