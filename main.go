package main

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"regexp"
	"strings"
)


func Find(cond string) (string, []string) {
	// Тут пишем условия на все операторы
	//fmt.Println(cond)
	re := regexp.MustCompile("((<=)|(>=)|(!=)|(LIKE)|=|<|>)")
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


func splitCondition(query string) squirrel.Sqlizer {
	exp, arr := Find(query)
	switch exp {
	case "=":
		// Где-то тут нужен колбэк
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


func Parse(query string, qb squirrel.SelectBuilder) (*squirrel.SelectBuilder, error){
	if checkAnotherOpers(query){
		panic(fmt.Errorf("unacceptable opertor"))
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
	fmt.Println(baseCond.ToSql())
	result := qb.Where(baseCond)
	return &result, nil
}


//func Parse(query string, qb squirrel.SelectBuilder) (*squirrel.SelectBuilder, error) {
//	// Распаристь строку
//	//
//
//	r := qb.Where(squirrel.And{squirrel.Eq{"test": 4}, squirrel.Or{squirrel.Eq{"test": 10}, squirrel.Gt{"test10": 20, "test123": 123}}})
//	fmt.Println(r.ToSql())
//	res := squirrel.Or{}
//	r1 :=  squirrel.Or{squirrel.Eq{"test": 10}, squirrel.Gt{"test10": 20, "test123": 123}}
//	r2 := squirrel.Or{squirrel.Eq{"wefwqe": 10}, squirrel.Gt{"wqer": 20}}
//	res = append(res, r1, r2)
//	fmt.Println(res.ToSql())
//	return &r, nil
//}

func main()  {
	qb := squirrel.SelectBuilder{}
	a, err := Parse("test = 10 OR test10 > 20 AND test123 > 123 OR wefwqe = 10 OR wqer > 20", qb)
	fmt.Println("Hello")
	fmt.Println(a, err)

}
