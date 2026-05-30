package ui

import (
	"context"
	"errors"
	"log"
	"strconv"
	"unicode/utf8"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/kakkky/hakoniwa/domain"
	"github.com/kakkky/hakoniwa/presentation/ui/components"
)

func (u *UI) registerResidentContainer(ctx context.Context) *dialog.CustomDialog {
	body := container.NewStack()

	navigate := func(content fyne.CanvasObject) {
		body.Objects = []fyne.CanvasObject{content}
		body.Refresh()
	}

	d := dialog.NewCustomWithoutButtons("", body, u.fyneMainWindow)
	d.Resize(fyne.NewSize(500, 300))

	// 初期状態としてformを表示させる
	navigate(u.registerResidentForm(ctx, navigate, func() { d.Dismiss() }))

	return d
}
func (u *UI) registerResidentForm(ctx context.Context, navigate func(fyne.CanvasObject), onClose func()) *fyne.Container {
	nameEntry := widget.NewEntry()
	ageEntry := widget.NewEntry()
	ageEntry.Validator = func(s string) error {
		if _, err := strconv.Atoi(s); err != nil {
			return errors.New("age must be number")
		}
		return nil
	}

	genderChoice := widget.NewRadioGroup([]string{"男", "女", "不明"}, func(string) {})

	personalityEntry := widget.NewMultiLineEntry()
	personalityEntry.Validator = func(s string) error {
		if utf8.RuneCountInString(s) > 100 {
			return errors.New("must be less than 100 chars")
		}
		return nil
	}
	personalityEntry.PlaceHolder = "性格や特徴を自由に設定してください(100文字以内)"

	var loading *widget.Activity
	onSubmitted := func() {
		navigate(components.LoadingContainer(loading))

		fyne.Do(func() {
			age, err := strconv.Atoi(ageEntry.Text)
			if err != nil {
			}

			var gender domain.Gender
			switch genderChoice.Selected {
			case "男":
				gender = domain.Male
			case "女":
				gender = domain.Female
			default:
				gender = domain.Unspecified
			}

			err = u.usecases.RegisterResident.Exec(ctx, nameEntry.Text, age, gender, personalityEntry.Text)
			if err != nil {
				log.Println(err)
			}
			navigate(u.registerResidentResult(ctx, navigate, onClose))
		})
	}

	return container.NewBorder(
		widget.NewLabelWithStyle("住人を登録しましょう", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		container.NewCenter(
			container.NewGridWithColumns(2,
				widget.NewButton("やめる", onClose),
				widget.NewButtonWithIcon("登録する", theme.ConfirmIcon(), onSubmitted),
			),
		),
		nil, nil,
		widget.NewForm(
			widget.NewFormItem("名前：", nameEntry),
			widget.NewFormItem("年齢：", ageEntry),
			widget.NewFormItem("性別", genderChoice),
			widget.NewFormItem("性格：", personalityEntry),
		),
	)
}

func (u *UI) registerResidentResult(ctx context.Context, navigate func(fyne.CanvasObject), onClose func()) *fyne.Container {
	return container.NewBorder(
		widget.NewLabelWithStyle("住人が登録されました", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		container.NewCenter(
			container.NewGridWithColumns(2,
				widget.NewButton("終了する", onClose),
				components.NewNavigateNextButton("続けて登録する", func() { navigate(u.registerResidentResult(ctx, navigate, onClose)) }),
			),
		),
		nil,
		nil,
		container.NewVBox(),
	)
}
