package utils

import (
	"bytes"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	//"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	//"net/http"
	"os/exec"
	"strconv"
)

//func ExeSh(trs system.SysAnaTrs, cc chan<- system.SysAnaTrs, r chan<- error) /*error*/ {
func ExeSh(trs system.SysAnaTrs, ws *websocket.Conn) error {
	//func ExeSh(c chan<- system.SysAnaTrs) error {
	//defer close(cc)
	ParaNumString := strconv.Itoa(trs.ParaNum)
	StrandString := ""
	if err := trs.Strand; err {
		StrandString = "yes"
	} else {
		StrandString = "no"
	}
	//           /work/run/project/PAMS2/bioinfo_pip/pip_mRNA.sh          /work/run/project/RNASeq/pip_mRNA.sh
	//cmd := exec.Command("bash", "/work/run/project/PAMS2/bioinfo_pip/pip_mRNA.sh", trs.ProjectNo, trs.RawCleanDir, trs.AnalysisType, trs.SampleInfo, trs.CmpGroup, trs.SpeciesNames,
	//	trs.DifferenceThreshold, trs.SubProjectNo, trs.PrjType, ParaNumString, trs.SeqType, StrandString, trs.GeneId, trs.FeatureType, trs.Category, trs.PpiSpecies, ">test.out", "2>test.err")
	defer ws.Close()
	var err error
	for {
		fmt.Println("webSocket 打通 前后端 传递信息 双向通道")
		ws.WriteMessage(1, []byte("分析中"))
		fmt.Println("开始执行shell脚本。。。")
		cmd := exec.Command("/bin/bash", "-c", "/work/run/project/PAMS2/bioinfo_pip/pip_mRNA.sh"+" "+trs.ProjectNo+" "+trs.RawCleanDir+" "+trs.AnalysisType+" "+trs.SampleInfo+" "+
			trs.CmpGroup+" "+trs.SpeciesNames+" "+trs.DifferenceThreshold+" "+trs.SubProjectNo+" "+trs.PrjType+" "+ParaNumString+" "+trs.SeqType+" "+StrandString+" "+trs.GeneId+" "+trs.FeatureType+" "+
			trs.Category+" "+trs.PpiSpecies+" "+">"+"/work/run/project/RNASeq/"+trs.ProjectNo+".out"+" "+"2>"+"/work/run/project/RNASeq/"+trs.ProjectNo+".err")
		var stdin, stdout, stderr bytes.Buffer
		cmd.Stdin = &stdin
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err = cmd.Run()
		//_ = cmd.Run()
		outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
		fmt.Printf("out:\n%s\nerr:\n%s\n", outStr, errStr)
		fmt.Println("shell脚本执行完毕。。。")
		err2 := ws.WriteMessage(1, []byte("分析完"))
		if err2 != nil {
			fmt.Println("webSocket write error")
			return err2
		} else {
			break
		}

	}
	/*fmt.Println("开始执行shell脚本。。。")
	cmd := exec.Command("/bin/bash", "-c", "/work/run/project/PAMS2/bioinfo_pip/pip_mRNA.sh"+" "+trs.ProjectNo+" "+trs.RawCleanDir+" "+trs.AnalysisType+" "+trs.SampleInfo+" "+
		trs.CmpGroup+" "+trs.SpeciesNames+" "+trs.DifferenceThreshold+" "+trs.SubProjectNo+" "+trs.PrjType+" "+ParaNumString+" "+trs.SeqType+" "+StrandString+" "+trs.GeneId+" "+trs.FeatureType+" "+
		trs.Category+" "+trs.PpiSpecies+" "+">"+"/work/run/project/RNASeq/"+trs.ProjectNo+".out"+" "+"2>"+"/work/run/project/RNASeq/"+trs.ProjectNo+".err")
	var stdin, stdout, stderr bytes.Buffer
	cmd.Stdin = &stdin
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	//err := cmd.Run()
	_ = cmd.Run()
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	fmt.Printf("out:\n%s\nerr:\n%s\n", outStr, errStr)
	fmt.Println("shell脚本执行完毕。。。")*/
	//cc <- trs
	//r <- err

	return err
}
