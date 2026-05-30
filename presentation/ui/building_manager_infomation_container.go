package ui

import (
	"context"
	"errors"
	"fmt"
	"image/color"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/kakkky/hakoniwa/domain"
	"github.com/kakkky/hakoniwa/presentation/ui/components"
)

func (u *UI) buildingManagerInfomationContainer(ctx context.Context) *fyne.Container {
	bdm, err := u.usecases.GetBuildingManager.Exec(ctx)
	if err != nil {

	}
	registrationBtn := withColor(
		color.NRGBA{R: 200, G: 200, B: 200, A: 255},
		container.NewPadded(
			container.NewCenter(
				u.registerBuildingManagerButton(
					func() { u.registerBuildingManagerFormDialog(ctx).Show() },
				),
			),
		),
	)

	buildingManagerInfomation := components.NewTappableContainer(
		container.NewHBox(
			widget.NewIcon(theme.AccountIcon()),
			widget.NewLabelWithStyle(
				fmt.Sprintf("管理人：%s", bdm.Name),
				fyne.TextAlignCenter,
				fyne.TextStyle{
					Bold: true,
				},
			),
		),
		func() {
			u.buildingManagerInfomationDialog(ctx).Show()
		},
	)

	stack := container.NewStack(
		registrationBtn,
		buildingManagerInfomation,
	)

	if bdm != nil {
		registrationBtn.Hide()
	} else {
		buildingManagerInfomation.Hide()
	}

	return stack
}

func (u *UI) registerBuildingManagerButton(onClick func()) *widget.Button {
	return widget.NewButtonWithIcon(
		"管理人登録",
		theme.ContentAddIcon(),
		onClick,
	)
}

func (u *UI) registerBuildingManagerFormDialog(ctx context.Context) *dialog.FormDialog {
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
		if err := u.usecases.RegisterBuildingManager.Exec(ctx, domain.NewBuildingManager(name, age, time.Now())); err != nil {

		}
		u.refreshTopBar(ctx)
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
		u.fyneMainWindow,
	)
}

func (u *UI) buildingManagerInfomationDialog(ctx context.Context) *dialog.CustomDialog {
	bdm, err := u.usecases.GetBuildingManager.Exec(ctx)
	if err != nil {

	}
	dialog := dialog.NewCustom(
		"管理人情報",
		"閉じる",
		container.NewVBox(
			components.NewEditableField("名前", string(bdm.Name), func(v string) {
				bdm.Name = domain.BuildManagerName(v)
				if err := u.usecases.UpdateBuildingManager.Exec(ctx, bdm); err != nil {

				}
				u.refreshTopBar(ctx)
			}),
			components.NewEditableField("年齢", strconv.Itoa(bdm.Age), func(v string) {
				age, err := strconv.Atoi(v)
				if err != nil {
					return
				}
				bdm.Age = age
				if err := u.usecases.UpdateBuildingManager.Exec(ctx, bdm); err != nil {

				}
				u.refreshTopBar(ctx)
			}),
			widget.NewLabel(fmt.Sprintf("就任：   %s", bdm.AppointedAt.Format("2006/01/02 15:04"))),
		),
		u.fyneMainWindow,
	)
	dialog.Resize(fyne.NewSize(500, 300))
	return dialog
}
