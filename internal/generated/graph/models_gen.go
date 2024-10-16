// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package graph

import (
	"fmt"
	"io"
	"strconv"
)

// Create User Input
type CreateUserInput struct {
	// 名前
	Name string `json:"name"`
	// メールアドレス
	Email string `json:"email"`
}

// Create User Output
type CreateUserOutput struct {
	// 作成した結果
	Status MutationStatus `json:"status"`
	// error message
	ErrorMessage *string `json:"errorMessage,omitempty"`
	// metadata
	Metadata *CreateUserOutputMetadata `json:"metadata,omitempty"`
}

type CreateUserOutputMetadata struct {
	// 作成したユーザ情報
	User *User `json:"user,omitempty"`
}

// Mutation
type Mutation struct {
}

// Query
type Query struct {
}

// User Information
type User struct {
	// ユーザID
	ID string `json:"id"`
	// 名前
	Name string `json:"name"`
	// メールアドレス
	Email string `json:"email"`
}

// Mutationの処理結果
type MutationStatus string

const (
	// 成功
	MutationStatusSuccess MutationStatus = "SUCCESS"
	// 作成済み
	MutationStatusAlreadyExists MutationStatus = "ALREADY_EXISTS"
	// 失敗
	MutationStatusFailure MutationStatus = "FAILURE"
	// バリデーションエラー
	MutationStatusValidationError MutationStatus = "VALIDATION_ERROR"
)

var AllMutationStatus = []MutationStatus{
	MutationStatusSuccess,
	MutationStatusAlreadyExists,
	MutationStatusFailure,
	MutationStatusValidationError,
}

func (e MutationStatus) IsValid() bool {
	switch e {
	case MutationStatusSuccess, MutationStatusAlreadyExists, MutationStatusFailure, MutationStatusValidationError:
		return true
	}
	return false
}

func (e MutationStatus) String() string {
	return string(e)
}

func (e *MutationStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = MutationStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid MutationStatus", str)
	}
	return nil
}

func (e MutationStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
