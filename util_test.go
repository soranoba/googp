package googp

import "testing"

func TestToSnake(t *testing.T) {
	assertEqual(t, toSnake("CreatedAt"), "created_at")
	assertEqual(t, toSnake("ID"), "id")
	assertEqual(t, toSnake("AbCdEf"), "ab_cd_ef")
}