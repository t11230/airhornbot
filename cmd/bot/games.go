package main

import (
    "bytes"
    "fmt"
    "text/tabwriter"
    "math/rand"
    "github.com/bwmarrin/discordgo"
    "strconv"
    "strings"
)

func rollDice(guild *discordgo.Guild, user *discordgo.User, args []string) string{
    draw := false
    r:=0
    maxnum:=0
    var err error
    if len(args)>1 {
        if strings.HasPrefix(args[1], "d") {
            maxnum, err = strconv.Atoi(strings.Replace(args[1], "d", "", 1))
            if err!=nil {
                return ""
            }
            if isValidDie(maxnum) {
                draw = true
            }
        } else {
            maxnum, err = strconv.Atoi(args[1])
            if err!=nil {
                return ""
            }
        }
        r = rand.Intn(maxnum) + 1
    } else {
        maxnum =6
        r = rand.Intn(6) + 1
        draw = true
    }
    w := &tabwriter.Writer{}
    buf := &bytes.Buffer{}
    result := ""
    if draw {
        if maxnum == 6 {
            C := "o "
            s:="---------\n| "+string(C[utilBooltoInt(r<=1)])+"   "+string(C[utilBooltoInt(r<=3)])+" |\n| "+string(C[utilBooltoInt(r<=5)])
            z:=string(C[utilBooltoInt(r<=5)])+" |\n| "+string(C[utilBooltoInt(r<=3)])+"   "+string(C[utilBooltoInt(r<=1)])+" |\n---------"
            result = s+" "+string(C[utilBooltoInt((r&1)==0)])+" "+z
        } else if (maxnum == 4) || (maxnum == 8) {
            result = "      *\n     * *\n    *   *\n   *  "+strconv.Itoa(r)+"  *\n  *       *\n * * * * * *"
        } else if maxnum == 10 {
            numstring := strconv.Itoa(r)
            if r > 9 {
                result = "        *\n       * *\n      *   *\n     * "+string(numstring[0])+" "+string(numstring[1])+" *\n      *   *\n        *"
            } else {
                result = "        *\n       * *\n      *   *\n     *  "+numstring+"  *\n      *   *\n        *"
            }
        } else if maxnum == 12 {
            numstring := strconv.Itoa(r)
            if r > 9 {
                result = "         *\n      *     *\n    *   "+string(numstring[0])+" "+string(numstring[1])+"   *\n     *       *\n      * * * *"
            } else {
                result = "         *\n      *     *\n    *    "+numstring+"    *\n     *       *\n      * * * *"
            }
        } else if maxnum == 20 {
            numstring := strconv.Itoa(r)
            if r > 9 {
                result = "      *\n     * *\n    *   *\n   * "+string(numstring[0])+" "+string(numstring[1])+" *\n  *       *\n * * * * * *"
            } else {
                result = "      *\n     * *\n    *   *\n   *  "+numstring+"  *\n  *       *\n * * * * * *"
            }
        }
    } else{
        result = "The result is: "+strconv.Itoa(r)
    }
    w.Init(buf, 0, 4, 0, ' ', 0)
    fmt.Fprintf(w, "```\n")
    fmt.Fprintf(w, result)
    fmt.Fprintf(w, "```\n")
    w.Flush()
    return buf.String()
}

func isValidDie(num int) bool {
    return utilIntInSlice(num, []int{4,6,8,10,12,20})
}
