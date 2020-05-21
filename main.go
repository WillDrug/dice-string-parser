package main

import (
	"flag"
	"fmt"
	"regexp"
	"math/rand"
	"sort"
	"strconv"
	"time"
	"sync"
)

var onlyOnce sync.Once
var dstring string
var re = regexp.MustCompile("^(?P<num>\\d+)(d(?P<dice>\\d+)((kh(?P<kh>\\d+))|(kl(?P<kl>\\d+))?)?)?(l(?P<lim>\\d+))?$")

type RollParse struct {
	dstring 	string
	arithmetic  []int // -2 div -1 sub 1 sum 2 mul ?
	rolls 		[]Roll
	total		int
}

type Roll struct{
	num  		int
	dice 		int
	keep 		int
	keep_high   bool
	lim 		int
	total		int	
	rolls 		[]int
	discarded   []int
}

func init() {
	flag.StringVar(&dstring, "d", "", "Dice string")
	flag.Parse()
}

func main() {
	var parsed RollParse
	var err error
	parsed, err = Parse(dstring)
	fmt.Println(parsed, err)	
}

func Parse(dstring string) (RollParse, error) {
	// 
	// example string with full functionality:
	// ((2d14 + 4d3kh1 - 2d8kl2 + 120)l300 * 10)l1000 + (something) - (rather)
	// example string for the current idea:
	// 1d20+5d10kh1l100+120
	var rp RollParse
	var err error
	rp = RollParse{dstring, make([]int, 0), make([]Roll, 0), 0}
	// start filling rolls
	var temp string
	var sym int
	// save roll result and the mathematical operation code
	for _, ch := range rp.dstring {
		switch ch {
			case '+':
				sym = 1
			case '-':
				sym = -1
			default:
				temp = temp + string(ch)
		}
		if sym != 0 {
			err = rp.Append(temp, sym)
			if err != nil {
				return RollParse{}, err
			}
			sym = 0
			temp = ""			
		}
	}
	// append the last roll
	err = rp.Append(temp, 0)
	if err != nil {
		return RollParse{}, err
	}
	// calculating total
	rp.total = rp.rolls[0].total
	for i, v := range rp.arithmetic {
		// for subtraction v is -1
		rp.total += v*rp.rolls[i+1].total
	}
	return rp, nil
}

func (rp *RollParse) Append(rstring string, arithmetic int) error {
	var roll Roll
	var err error
	roll, err = Rollify(rstring)
	if err != nil {
		return fmt.Errorf("Failed to parse %s somehow: %s", rstring, err)
	}
	rp.rolls = append(rp.rolls, roll)
	if arithmetic != 0 {
		rp.arithmetic = append(rp.arithmetic, arithmetic)
	}
	return nil
}

func Rollify(rstring string) (Roll, error) {
	// take a base string containing a roll, validate it and turn into Roll struct
	// validate
	if !Validate(rstring) {
		return Roll{}, fmt.Errorf("Error: %s is not a valid dice string", rstring)
	}
	// from this point onwards we think the roll string is valid
	// match and save named capture groups
	match := re.FindStringSubmatch(rstring)
    result := make(map[string]string)
    for i, name := range re.SubexpNames() {
        if i != 0 && name != "" {
            result[name] = match[i]
        }
    }
    // crate struct
    var roll Roll
    roll = Roll{0, 1, 0, true, 0, 0, nil, nil}
    // turn capture groups into struct fields
    // num, dice, keep, keep_high, lim, total
    var nv int64
    var err error
    for k, v := range result{
    	if (v == "") {
    		continue
    	}
    	nv, err = strconv.ParseInt(v, 10, 64)
    	if err != nil {
    		return Roll{}, fmt.Errorf("Found %s in %s which is not an integer", v, k)
    	}
    	switch k {
	    	case "kl":
	    		roll.keep_high = true
	    		fallthrough
	    	case "kh":
	    		roll.keep = int(nv) // bad practices ahoy
	    	case "num":
	    		roll.num = int(nv)
	    	case "dice":
	    		roll.dice = int(nv)
	    	case "lim":
	    		roll.lim = int(nv)
	    }
    }
    // todo: move all the separate dice functions to code functions
    // calculate total
    // skip rolling function if dice are 1-sided
    if (roll.dice == 1) {    	
		roll.total = roll.num
    } else {
    	roll.rolls = make([]int, roll.num)
    	for i:=0;i<roll.num;i++ {
	    	roll.rolls[i] = RollDice(roll.dice)
	    }
	}
    if (roll.keep > 0 && roll.dice > 1) {
    	if (roll.keep > roll.num) { // we are checking this only if we are keeping, because 0 is the default non-keep roll
    		return Roll{}, fmt.Errorf("Keep value is higher than the number of dice")
    	}
    	if (roll.keep < roll.num) { // if keep == num, the result stays the same
    		roll.discarded = make([]int, roll.num-roll.keep)
    		sort.Ints(roll.rolls)
    		if !roll.keep_high {
    			sort.Reverse(sort.IntSlice(roll.rolls))
    		} 
    		for i:=0;i<roll.num-roll.keep;i++{
    			roll.discarded[i] = roll.rolls[i]
    		}
    		roll.rolls = roll.rolls[roll.num-roll.keep:]
    	}
    }
    for _, v := range roll.rolls {
    	roll.total += v
    }
    if (roll.total > roll.lim && roll.lim > 0) {
    		roll.total = roll.lim
	}
    return roll, nil
}

func RollDice(sides int) int {
	// edge case
	if (sides == 1) {
		return 1
	}
	// randomizing 
	onlyOnce.Do(func() {
 		rand.Seed(time.Now().UnixNano()) // only run once
 	})
	// todo: add some fairness
	return rand.Intn(sides)+1
}

func Validate(rstring string) bool {
	return re.Match([]byte(rstring))
}

// ROLL l 1000
// ROLL l 1300 * ROLL (10)
// ROLL + ROLL - ROLL + ROLL(120)
// 