package lib

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/code"
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/image"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"github.com/hazemKrimi/crimson-vault/internal/types"
)

func GetConfigDirectory() (string, error) {
	home, err := os.UserHomeDir()

	if err != nil {
		return "", err
	}

	config, err := filepath.Abs(filepath.Join(home, DEFAULT_CONFIG_DIRECTORY))

	return config, nil
}

func SaveSession(session *sessions.Session, context echo.Context) error {
	if err := session.Save(context.Request(), context.Response()); err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: err, Messages: []string{"Unexpected error saving User session!"}}
	}

	return nil
}

func CreateSession(session *sessions.Session, context echo.Context, user *types.User) error {
	if err := uuid.Validate(user.SessionID); err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: err, Messages: []string{"Unexpected error saving User session!"}}
	}

	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
	}
	session.Values["id"] = user.ID
	session.Values["sessionId"] = user.SessionID
	session.Values["username"] = user.Username

	if err := SaveSession(session, context); err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: err, Messages: []string{"Unexpected error saving User session!"}}
	}

	return nil
}

func DeleteSession(session *sessions.Session, context echo.Context) error {
	session.Options.MaxAge = -1

	if err := SaveSession(session, context); err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: err, Messages: []string{"Unexpected error saving User session!"}}
	}

	return nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func getGrayColor() *props.Color {
	return &props.Color{
		Red:   200,
		Green: 200,
		Blue:  200,
	}
}

func getDarkGrayColor() *props.Color {
	return &props.Color{
		Red:   200,
		Green: 200,
		Blue:  200,
	}
}

func getBlueColor() *props.Color {
	return &props.Color{
		Red:   10,
		Green: 10,
		Blue:  150,
	}
}

func getRedColor() *props.Color {
	return &props.Color{
		Red:   150,
		Green: 10,
		Blue:  10,
	}
}

func getPageHeader(client types.Client, logo string) core.Row {
	return row.New(20).Add(
		image.NewFromFileCol(3, logo, props.Rect{
			Center:  true,
			Percent: 80,
		}),
		col.New(6),
		col.New(3).Add(
			text.New(fmt.Sprintf("%s, %s, %s", client.Address, client.Zip, client.Country), props.Text{
				Size:  8,
				Align: align.Right,
				Color: getRedColor(),
			}),
			text.New(client.Phone, props.Text{
				Top:   12,
				Style: fontstyle.BoldItalic,
				Size:  8,
				Align: align.Right,
				Color: getBlueColor(),
			}),
		),
	)
}

func getPageFooter(user types.User) core.Row {
	return row.New(20).Add(
		col.New(12).Add(
			text.New(user.Phone, props.Text{
				Top:   13,
				Style: fontstyle.BoldItalic,
				Size:  8,
				Align: align.Left,
				Color: getBlueColor(),
			}),
		),
	)
}

func getTransactions(items []types.Item) []core.Row {
	rows := []core.Row{
		row.New(5).Add(
			col.New(3),
			text.NewCol(4, "Service", props.Text{Size: 9, Align: align.Center, Style: fontstyle.Bold}),
			text.NewCol(2, "Quantity", props.Text{Size: 9, Align: align.Center, Style: fontstyle.Bold}),
			text.NewCol(3, "Price", props.Text{Size: 9, Align: align.Center, Style: fontstyle.Bold}),
		),
	}

	var contentsRow []core.Row
	for i, item := range items {
		r := row.New(4).Add(
			col.New(3),
			text.NewCol(4, item.Name, props.Text{Size: 8, Align: align.Center}),
			text.NewCol(2, fmt.Sprintf("%d", item.Quantity), props.Text{Size: 8, Align: align.Center}),
			text.NewCol(3, fmt.Sprintf("%d", item.Price), props.Text{Size: 8, Align: align.Center}),
		)
		if i%2 == 0 {
			gray := getGrayColor()
			r.WithStyle(&props.Cell{BackgroundColor: gray})
		}

		contentsRow = append(contentsRow, r)
	}

	rows = append(rows, contentsRow...)

	rows = append(rows, row.New(20).Add(
		col.New(7),
		text.NewCol(2, "Total:", props.Text{
			Top:   5,
			Style: fontstyle.Bold,
			Size:  8,
			Align: align.Right,
		}),
		text.NewCol(3, "$2.567,00", props.Text{
			Top:   5,
			Style: fontstyle.Bold,
			Size:  8,
			Align: align.Center,
		}),
	))

	return rows
}

func GenerateInvoice(invoice types.Invoice, user types.User, client types.Client) error {
	cfg := config.NewBuilder().WithPageNumber().WithLeftMargin(10).WithRightMargin(10).WithTopMargin(15).Build()
	darkGray := getDarkGrayColor()
	mrt := maroto.New(cfg)
	m := maroto.NewMetricsDecorator(mrt)

	err := m.RegisterHeader(getPageHeader(client, user.Logo))

	if err != nil {
		return err
	}

	err = m.RegisterFooter(getPageFooter(user))

	if err != nil {
		return err
	}

	m.AddRows(text.NewRow(10, "Invoice ABC123456789", props.Text{
		Top:   3,
		Style: fontstyle.Bold,
		Align: align.Center,
	}))

	m.AddRow(7,
		text.NewCol(3, "Transactions", props.Text{
			Top:   1.5,
			Size:  9,
			Style: fontstyle.Bold,
			Align: align.Center,
			Color: &props.WhiteColor,
		}),
	).WithStyle(&props.Cell{BackgroundColor: darkGray})

	m.AddRows(getTransactions(invoice.Items)...)

	m.AddRow(15,
		col.New(6).Add(
			code.NewBar("5123.151231.512314.1251251.123215", props.Barcode{
				Percent: 0,
				Proportion: props.Proportion{
					Width:  20,
					Height: 2,
				},
			}),
			text.New("5123.151231.512314.1251251.123215", props.Text{
				Top:    12,
				Family: "",
				Style:  fontstyle.Bold,
				Size:   9,
				Align:  align.Center,
			}),
		),
		col.New(6),
	)

	document, err := m.Generate()

	if err != nil {
		return err
	}

	err = document.Save(fmt.Sprintf("invoices/%s.pdf", invoice.ID))

	if err != nil {
		return err
	}

	err = document.GetReport().Save(fmt.Sprintf("invoices/%s.txt", invoice.ID))

	if err != nil {
		return err
	}

	return nil
}
