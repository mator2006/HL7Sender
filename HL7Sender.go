package main

import (
	"bufio"
	"fmt"
	"gopkg.in/gcfg.v1"
	"net"
	"os"
	"strings"
)

type hl7 struct {
	Host             string
	Port             string
	hl7mesg          string
	hl7mesgsub       string
	PatientNAME      string
	PatientBD        string
	PatientSEX       string
	PatientID        string
	OrderControl     string
	AccessionNO      string
	RPNO             string
	SPNO             string
	MODALITY         string
	SPSDESC          string
	StationNAME      string
	StationAET       string
	StudyInstanceUID string
}

func main() {

	configurefile := "config.ini"
	hl7data := hl7{}
	iniexist, err := os.Stat(configurefile)

	if iniexist != nil {

		fmt.Println("Found Configure file.\nConfigure file is [" + configurefile + "].")
		config := struct {
			Netconn struct {
				Host string
				Port string
			}
			HL7order struct {
				PatientID        string
				PatientNAME      string
				PatientBD        string
				PatientSEX       string
				OrderControl     string
				AccessionNO      string
				RPNO             string
				SPNO             string
				MODALITY         string
				SPSDESC          string
				StationNAME      string
				StationAET       string
				StudyInstanceUID string
			}
		}{}

		err := gcfg.ReadFileInto(&config, configurefile)
		if err != nil {
			fmt.Println(err, "Failed to parse[fy:fen xi] configure file:\nYou can DO THIS:\n1.Tyr modify the configure file:["+configurefile+"]\n2.Or remove this configure files,Will use DEFAULT configure")
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
			hl7data.hl7mesgsub = `MSH|^~\&|||||||ORM^O01||P|2.3|||||CN` + "\r" + `PID|||` + hl7data.PatientID + `^^^HIS||` + hl7data.PatientNAME + `||` + hl7data.PatientBD + `|` + hl7data.PatientSEX + `|||||||||||||||||||||||||||||||` + "\r" +
				`PV1||E|||||||||||||||||V103-1^^^ADT1||||||||||||||||||||||||||||||||V|` + "\r" +
				`ORC|` + hl7data.OrderControl + `||||SC||1^once^^^^S||||||||||||` + "\r" +
				`OBR|1|||^^^^` + hl7data.SPSDESC + `||||||||||||||` + hl7data.AccessionNO + `|` + hl7data.RPNO + `|` + hl7data.SPNO + `||||` + hl7data.MODALITY +
				`|||1^once^^^^S|||WALK||||||||||||||P1^` + hl7data.SPSDESC + `^ERL_MESA, ZDS|1.113654.3.13.1025^100^Application^DICOM` + "\r" +
				`ZDS|` + hl7data.StudyInstanceUID + `^MEPACS^Application^DICOM|` + hl7data.StationNAME + `|` + hl7data.StationAET + "\r"
			hl7data.hl7mesg = "\v" + hl7data.hl7mesgsub + "\x1C" + "\r" //谨慎同上

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
	} else {
		fmt.Println(err, "Error.")
		return
	}
}
