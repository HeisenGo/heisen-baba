package presenter

import (
	"tripcompanyservice/internal/invoice"
	"tripcompanyservice/pkg/fp"
)

type AgencyInvoice struct {
	ID             uint         `json:"id"`
	TicketID       uint         `json:"ticket_id"`
	Ticket         AgencyTicket `json:"ticket"`
	IssuedDate     Timestamp    `json:"issued_date"`
	Info           string       `json:"info"`
	PerAmountPrice float64      `json:"per_price"`
	TotalPrice     float64      `json:"total_price"`
}

type UserInvoice struct {
	ID             uint       `json:"id"`
	TicketID       uint       `json:"ticket_id"`
	Ticket         UserTicket `json:"ticket"`
	IssuedDate     Timestamp  `json:"issued_date"`
	Info           string     `json:"info"`
	PerAmountPrice float64    `json:"per_price"`
	TotalPrice     float64    `json:"total_price"`
}

func InvoiceToAgencyInvoice(i invoice.Invoice) AgencyInvoice {
	ticket := TicketToAgencyTicket(*i.Ticket)
	return AgencyInvoice{
		ID:             i.ID,
		TicketID:       i.TicketID,
		Ticket:         ticket,
		IssuedDate:     Timestamp(i.IssuedDate),
		Info:           i.Info,
		PerAmountPrice: i.PerAmountPrice,
		TotalPrice:     i.TotalPrice,
	}
}

func InvoiceToUserInvoice(i invoice.Invoice) UserInvoice {
	ticket := TicketToUserTicket(*i.Ticket)
	return UserInvoice{
		ID:             i.ID,
		TicketID:       i.TicketID,
		Ticket:         ticket,
		IssuedDate:     Timestamp(i.IssuedDate),
		Info:           i.Info,
		PerAmountPrice: i.PerAmountPrice,
		TotalPrice:     i.TotalPrice,
	}
}

func BatchInvoicesToUserInvoices(is []invoice.Invoice) []UserInvoice {
	return fp.Map(is, InvoiceToUserInvoice)
}

func BatchInvoicesToAgencyInvoices(is []invoice.Invoice) []AgencyInvoice {
	return fp.Map(is, InvoiceToAgencyInvoice)
}
