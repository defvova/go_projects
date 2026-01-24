package handler

import (
	"calculator/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Response struct {
	Num1   int    `json:"num1"`
	Num2   int    `json:"num2"`
	Opt    string `json:"opt"`
	Result int    `json:"result"`
}

func HandleAdd(w http.ResponseWriter, r *http.Request) {
	if ok := setMethod(w, r, "POST"); !ok {
		return
	}
	num1, num1Ok := parseNum(w, r, "num1")
	num2, num2Ok := parseNum(w, r, "num2")
	if !num1Ok || !num2Ok {
		return
	}

	o := service.Operation{Num1: num1, Num2: num2}
	body := Response{Num1: num1, Num2: num2, Opt: "add", Result: o.Add()}

	json.NewEncoder(w).Encode(body)
}

func HandleSubtract(w http.ResponseWriter, r *http.Request) {
	if ok := setMethod(w, r, "POST"); !ok {
		return
	}
	num1, num1Ok := parseNum(w, r, "num1")
	num2, num2Ok := parseNum(w, r, "num2")
	if !num1Ok || !num2Ok {
		return
	}

	o := service.Operation{Num1: num1, Num2: num2}
	body := Response{Num1: num1, Num2: num2, Opt: "subtract", Result: o.Subtract()}

	json.NewEncoder(w).Encode(body)
}

func HandleMultiply(w http.ResponseWriter, r *http.Request) {
	if ok := setMethod(w, r, "POST"); !ok {
		return
	}
	num1, num1Ok := parseNum(w, r, "num1")
	num2, num2Ok := parseNum(w, r, "num2")
	if !num1Ok || !num2Ok {
		return
	}

	o := service.Operation{Num1: num1, Num2: num2}
	body := Response{Num1: num1, Num2: num2, Opt: "multiply", Result: o.Multiply()}

	json.NewEncoder(w).Encode(body)
}

func HandleDivide(w http.ResponseWriter, r *http.Request) {
	if ok := setMethod(w, r, "POST"); !ok {
		return
	}
	num1, num1Ok := parseNum(w, r, "num1")
	num2, num2Ok := parseNum(w, r, "num2")
	if !num1Ok || !num2Ok {
		return
	}

	o := service.Operation{Num1: num1, Num2: num2}
	body := Response{Num1: num1, Num2: num2, Opt: "divide", Result: o.Divide()}

	json.NewEncoder(w).Encode(body)
}

type RequestParam struct {
	Items []int `json:"items,omitempty"`
}

func HandleSum(w http.ResponseWriter, r *http.Request) {
	if ok := setMethod(w, r, "POST"); !ok {
		return
	}
	var req RequestParam
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	o := service.Operation{}
	body := Response{Num1: 0, Num2: 0, Opt: "sum", Result: o.Sum(req.Items)}
	json.NewEncoder(w).Encode(body)
}

func setMethod(w http.ResponseWriter, r *http.Request, m string) bool {
	if r.Method != m {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return false
	}
	w.Header().Set("Content-Type", "application/json")
	return true
}

func parseNum(w http.ResponseWriter, r *http.Request, name string) (int, bool) {
	val := r.FormValue(name)
	num, err := strconv.Atoi(val)

	if err != nil {
		msg := fmt.Sprintf("`%v` can`t be blank and must include only int value.", name)
		http.Error(w, msg, http.StatusBadRequest)
		return 0, false
	}

	return num, true
}
