package components

import (
	"context"
	"errors"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kakkky/hakoniwa/domain"
	"github.com/kakkky/hakoniwa/usecase"
)

func RegisterBuildingManagerFormDialog(ctx context.Context, registerBuildingManager *usecase.RegisterBuildingManager, parentWindow fyne.Window, refresh func()) *dialog.FormDialog {
	nameEntry := widget.NewEntry()
	ageEntry := widget.NewEntry()
	ageEntry.Validator = func(s string) error {
		if _, err := strconv.Atoi(s); err != nil {
			return errors.New("age must be number")
		}
		return nil
	}

	onSubmit := func(isConfirm bool) {
		if !isConfirm {
			return
		}
		name := nameEntry.Text
		age, err := strconv.Atoi(ageEntry.Text)
		if err != nil {

		}
		if err := registerBuildingManager.Exec(ctx, domain.NewBuildingManager(name, age, time.Now())); err != nil {

		}
		refresh()
	}

	return dialog.NewForm(
		"あなたをマンションの管理人として登録します",
		"登録",
		"やめる",
		[]*widget.FormItem{
			widget.NewFormItem("名前", nameEntry),
			widget.NewFormItem("年齢", ageEntry),
		},
		onSubmit,
		parentWindow,
	)
}
