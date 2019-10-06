package parser

import (
    "fmt"
    "text/scanner"
    "strings"
    "os"
    "bufio"
    "log"
)

type ScannerError struct {
    ErrorString string
    pos scanner.Position    
    tokenText string
}

func (err *ScannerError) Error() string { return "error at " + err.pos.String() +  err.ErrorString }

var (
    ErrBadToken     = "bad token."
    ErrBadSyntax    = "bad syntax."
    ErrBadInput     = "input is empty."
    ErrNotSupported = "not supported."
)

func digitaval(ch byte) int {
    switch {
    case '0' <= ch && ch <= '9':
        return int(ch - '0')
    }
    return 16
}


func scanIpOctet(p *scanner.Scanner, nth int) (uint8, *ScannerError) {
    var oct = make([]byte, 3) 
    var tok rune
    var base = []int{2,5,5}
    var l int
    var v int = 0

    tok = p.Scan()
    if (tok == scanner.Int) {
        l = copy(oct, p.TokenText());
        i := 0
                fmt.Printf("len oct %d\n", l)
        for (i < l) {
            dv := digitaval(oct[i]) 
            if (dv > base[i]) {
                break
             }
                v = v * 10 + dv
            i++
        } 

        if (i != l) {
            fmt.Printf("error/ octet is %v\n", oct)
            return 0, &ScannerError{ErrorString: ErrBadToken,
                        pos: p.Pos(),
                        tokenText: p.TokenText()} 
        }

        if (nth != 3) {
            tok = p.Scan()
            if (tok != '.') {
                    return 0, &ScannerError{ErrorString: ErrBadSyntax,
                        pos: p.Pos(),
                        tokenText: p.TokenText()}
            }
            fmt.Println("-- At position", p.Pos(), " text:", p.TokenText())
            fmt.Println("-- token :", tok) 
        } 
        return uint8(v), nil
    } else {
             fmt.Printf("error/ tok is %d, tokenText is %v\n", tok, p.TokenText())
             fmt.Println("Err -- At position", p.Pos(), " text:", p.TokenText())
             return 0, &ScannerError{ErrorString: ErrBadSyntax,
                        pos: p.Pos(),
                        tokenText: p.TokenText()}
    }
    return 0, &ScannerError{ErrorString: ErrBadInput,
                                                pos: p.Pos(),
                                                tokenText: p.TokenText()}
}

func scanIpAddr(p scanner.Scanner) {

    var oct1, oct2, oct3, oct4 uint8
    var err *ScannerError

    oct1, err =  scanIpOctet(&p, 0)
    fmt.Printf("octet1 is %v\n", oct1)
    fmt.Printf("err is %v\n", err)
    if (err == nil) {
        oct2, err =  scanIpOctet(&p, 1)
        fmt.Printf("octet2 is %v\n", oct2)
        fmt.Printf("err is %v\n", err)
    }
    if (err == nil) {
        oct3, err =  scanIpOctet(&p, 2)
        fmt.Printf("octet3 is %v\n", oct3)
        fmt.Printf("err is %v\n", err)
    }
    if (err == nil) {
        oct4, err =  scanIpOctet(&p, 3)
        fmt.Printf("octet4 is %v\n", oct4)
        fmt.Printf("err is %v\n", err)
    }

    fmt.Println("ipaddr", oct1, ".", oct2, ".", oct3, ".", oct4)

}

func init() {
}

func Tokenize(line string) {
    // Initialize the scanner.
    var s scanner.Scanner
    s.Filename = "file.text"
    s.Init(strings.NewReader(line))
    s.Mode = (scanner.GoTokens & ^scanner.ScanFloats)
    s.Mode |= scanner.ScanInts

    // Repeated calls to Scan yield the token sequence found in the input.
    for {
        tok := s.Scan()
        if tok == scanner.EOF {
            break
        }
        fmt.Printf("%s\t%s\t%s\n", s.Pos(), s.TokenText(), scanner.TokenString(tok))
    }
}

func ParseFile(filename string) {
    file, err := os.Open(filename)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)

    // Default scanner is bufio.ScanLines. Lets use ScanWords.
    // Could also use a custom function of SplitFunc type
    // scanner.Split(bufio.ScanWords)

    // Scan for next token.
    success := true
    for success {
        success = scanner.Scan()
        if (success) {
            fmt.Println("line = " + scanner.Text())  // token in unicode-char
            Tokenize(scanner.Text())
        }
        // fmt.Println(scanner.Bytes()) // token in bytes
    }

    if success == false {
        // False on error or EOF. Check error
        err = scanner.Err()
        if err == nil {
            log.Println("Scan completed and reached EOF")
        } else {
            log.Fatal(err)
        }
    }
}
