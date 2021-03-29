package main

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"regexp"
	"strings"
	//"reflect"
)


func Find(cond string) (string, []string) {
	re := regexp.MustCompile("((<=)|(>=)|(!=)|(!~)|(~)|=|<|>)")
	found := re.FindAllString(cond, 1)
	if len(found) < 1{
		panic(fmt.Errorf("has not cond operator"))
	}
	CondArr := strings.Split(cond, found[0])
	if len(CondArr) > 1 {
		return found[0], CondArr
	}
	return "", []string{}
}


func splitCondition(condition string) squirrel.Sqlizer {
	exp, arr := Find(condition)
	//if !callBackHandler(arr[0], arr[1]){
	//panic(fmt.Errorf("type mismatch"))
	//}
	switch exp {
	case "=":
		return squirrel.Eq{arr[0]: arr[1]}
	case "<=":
		return squirrel.LtOrEq{arr[0]: arr[1]}
	case ">=":
		return squirrel.GtOrEq{arr[0]: arr[1]}
	case "!=":
		return squirrel.NotEq{arr[0]: arr[1]}
	case ">":
		return squirrel.Gt{arr[0]: arr[1]}
	case "<":
		return squirrel.Lt{arr[0]: arr[1]}
	case "~":
		return squirrel.Like{arr[0]: arr[1]}
	case "!~":
		return squirrel.NotLike{arr[0]: arr[1]}
	}
	return nil
}


func splitAndCodition(andCond []string)  squirrel.Sqlizer{
	resAndCond := squirrel.And{}
	for _, cond := range andCond {
		resAndCond = append(resAndCond, splitCondition(cond))
	}
	return resAndCond
}


func checkAnotherOpers(query string) (flag bool) {
	anOpers := []string{"LIMIT", "ORDER BY", "GROP BY", "HAVING", "DISTINCT", "JOIN"}
	for _, oper := range anOpers{
		flag = strings.Contains(query, oper)
	}
	return
}

//func callBackHandler(name string, value string) bool{
//	partsName := strings.Split(name,".")
//	//nType = discover type of colomn partsName(len(partsName)-1)
//	//compare the value type and the resulting name type
//	if reflect.TypeOf(value) == nType{
//		return true
//	} else{
//		return false
//	}
//}

func Parse(query string, qb squirrel.SelectBuilder) (*squirrel.SelectBuilder, error){
	if checkAnotherOpers(query){
		return nil, fmt.Errorf("unacceptable opertor")
	}
	conditions := strings.Split(query, "OR")
	baseCond := squirrel.Or{}
	for _, oneOrCond := range conditions {
		andConditions := strings.Split(oneOrCond,"AND")
		if len(andConditions) > 1 {
			baseCond = append(baseCond, splitAndCodition(andConditions))
		} else {
			baseCond = append(baseCond, splitCondition(oneOrCond))
		}
	}
	//fmt.Println(baseCond.ToSql())
	result := qb.Where(baseCond)
	return &result, nil
}


func main()  {
	qb := squirrel.SelectBuilder{}
	a, err := Parse("test = 10 OR test10 > 20 AND test123 > 123 OR wefwqe = 10 OR wqer > 20", qb)
	if err != nil {
		panic(fmt.Errorf("some error in parsing"))
	}
	fmt.Println(a)
}
