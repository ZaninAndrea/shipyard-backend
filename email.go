package main

import (
	"bytes"
	"os"
	"strconv"
	"text/template"
	"time"

	gomail "gopkg.in/gomail.v2"
)

const rawEmailTemplate string = `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional //EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
    <html
      xmlns="http://www.w3.org/1999/xhtml"
      xmlns:o="urn:schemas-microsoft-com:office:office"
      xmlns:v="urn:schemas-microsoft-com:vml"
    >
      <head>
        <!--[if gte mso 9
          ]><xml
            ><o:OfficeDocumentSettings
              ><o:AllowPNG /><o:PixelsPerInch
                >96</o:PixelsPerInch
              ></o:OfficeDocumentSettings
            ></xml
          ><!
        [endif]-->
        <meta content="text/html; charset=utf-8" http-equiv="Content-Type" />
        <meta content="width=device-width" name="viewport" />
        <!--[if !mso]><!-->
        <meta content="IE=edge" http-equiv="X-UA-Compatible" />
        <!--<![endif]-->
        <title></title>
        <!--[if !mso]><!-->
        <link
          href="https://fonts.googleapis.com/css?family=Roboto"
          rel="stylesheet"
          type="text/css"
        />
        <!--<![endif]-->
        <style type="text/css">
          body {
            margin: 0;
            padding: 0;
          }
    
          table,
          td,
          tr {
            vertical-align: top;
            border-collapse: collapse;
          }
    
          * {
            line-height: inherit;
          }
    
          a[x-apple-data-detectors="true"] {
            color: inherit !important;
            text-decoration: none !important;
          }
        </style>
        <style id="media-query" type="text/css">
          @media (max-width: 670px) {
            .block-grid,
            .col {
              min-width: 320px !important;
              max-width: 100% !important;
              display: block !important;
            }
    
            .block-grid {
              width: 100% !important;
            }
    
            .col {
              width: 100% !important;
            }
    
            .col > div {
              margin: 0 auto;
            }
    
            img.fullwidth,
            img.fullwidthOnMobile {
              max-width: 100% !important;
            }
    
            .no-stack .col {
              min-width: 0 !important;
              display: table-cell !important;
            }
    
            .no-stack.two-up .col {
              width: 50% !important;
            }
    
            .no-stack .col.num4 {
              width: 33% !important;
            }
    
            .no-stack .col.num8 {
              width: 66% !important;
            }
    
            .no-stack .col.num4 {
              width: 33% !important;
            }
    
            .no-stack .col.num3 {
              width: 25% !important;
            }
    
            .no-stack .col.num6 {
              width: 50% !important;
            }
    
            .no-stack .col.num9 {
              width: 75% !important;
            }
    
            .video-block {
              max-width: none !important;
            }
    
            .mobile_hide {
              min-height: 0px;
              max-height: 0px;
              max-width: 0px;
              display: none;
              overflow: hidden;
              font-size: 0px;
            }
    
            .desktop_hide {
              display: block !important;
              max-height: none !important;
            }
          }
        </style>
      </head>
      <body
        class="clean-body"
        style="
          margin: 0;
          padding: 0;
          -webkit-text-size-adjust: 100%;
          background-color: #ffffff;
        "
      >
        <!--[if IE]><div class="ie-browser"><![endif]-->
        <table
          bgcolor="#FFFFFF"
          cellpadding="0"
          cellspacing="0"
          class="nl-container"
          role="presentation"
          style="
            table-layout: fixed;
            vertical-align: top;
            min-width: 320px;
            margin: 0 auto;
            border-spacing: 0;
            border-collapse: collapse;
            mso-table-lspace: 0pt;
            mso-table-rspace: 0pt;
            background-color: #ffffff;
            width: 100%;
          "
          valign="top"
          width="100%"
        >
          <tbody>
            <tr style="vertical-align: top;" valign="top">
              <td style="word-break: break-word; vertical-align: top;" valign="top">
                <!--[if (mso)|(IE)]><table width="100%" cellpadding="0" cellspacing="0" border="0"><tr><td align="center" style="background-color:#FFFFFF"><![endif]-->
                <div style="background-color: {{.HeaderColor}};">
                  <div
                    class="block-grid"
                    style="
                      margin: 0 auto;
                      min-width: 320px;
                      max-width: 650px;
                      overflow-wrap: break-word;
                      word-wrap: break-word;
                      word-break: break-word;
                      background-color: transparent;
                    "
                  >
                    <div
                      style="
                        border-collapse: collapse;
                        display: table;
                        width: 100%;
                        background-color: transparent;
                      "
                    >
                      <!--[if (mso)|(IE)]><table width="100%" cellpadding="0" cellspacing="0" border="0" style="background-color:{{.HeaderColor}};"><tr><td align="center"><table cellpadding="0" cellspacing="0" border="0" style="width:650px"><tr class="layout-full-width" style="background-color:transparent"><![endif]-->
                      <!--[if (mso)|(IE)]><td align="center" width="650" style="background-color:transparent;width:650px; border-top: 0px solid transparent; border-left: 0px solid transparent; border-bottom: 0px solid transparent; border-right: 0px solid transparent;" valign="top"><table width="100%" cellpadding="0" cellspacing="0" border="0"><tr><td style="padding-right: 0px; padding-left: 0px; padding-top:30px; padding-bottom:30px;"><![endif]-->
                      <div
                        class="col num12"
                        style="
                          min-width: 320px;
                          max-width: 650px;
                          display: table-cell;
                          vertical-align: top;
                          width: 650px;
                        "
                      >
                        <div style="width: 100% !important;">
                          <!--[if (!mso)&(!IE)]><!-->
                          <div
                            style="
                              border-top: 0px solid transparent;
                              border-left: 0px solid transparent;
                              border-bottom: 0px solid transparent;
                              border-right: 0px solid transparent;
                              padding-top: 30px;
                              padding-bottom: 30px;
                              padding-right: 0px;
                              padding-left: 0px;
                            "
                          >
                            <!--<![endif]-->
                            <div
                              align="center"
                              class="img-container center fixedwidth"
                              style="padding-right: 0px; padding-left: 0px;"
                            >
                              <!--[if mso]><table width="100%" cellpadding="0" cellspacing="0" border="0"><tr style="line-height:0px"><td style="padding-right: 0px;padding-left: 0px;" align="center"><!
                              [endif]--><a
                                href="{{.Domain}}"
                                style="outline: none;"
                                tabindex="-1"
                                target="_blank"
                              >
                                <img
                                  align="center"
                                  alt="logo"
                                  border="0"
                                  class="center fixedwidth"
                                  src="{{.LogoLink}}"
                                  style="
                                    text-decoration: none;
                                    -ms-interpolation-mode: bicubic;
                                    border: 0;
                                    width: auto;
									margin: auto;
                                    height: 128px;
                                    display: block;
                                  "
                                  title="logo"
                                  height="128"
                              /></a>
                              <!--[if mso]></td></tr></table><![endif]-->
                            </div>
                            <!--[if (!mso)&(!IE)]><!-->
                          </div>
                          <!--<![endif]-->
                        </div>
                      </div>
                      <!--[if (mso)|(IE)]></td></tr></table><![endif]-->
                      <!--[if (mso)|(IE)]></td></tr></table></td></tr></table><![endif]-->
                    </div>
                  </div>
                </div>
                <div style="background-color: #fff;">
                  <div
                    class="block-grid"
                    style="
                      margin: 0 auto;
                      min-width: 320px;
                      max-width: 650px;
                      overflow-wrap: break-word;
                      word-wrap: break-word;
                      word-break: break-word;
                      background-color: transparent;
                    "
                  >
                    <div
                      style="
                        border-collapse: collapse;
                        display: table;
                        width: 100%;
                        background-color: transparent;
                      "
                    >
                      <!--[if (mso)|(IE)]><table width="100%" cellpadding="0" cellspacing="0" border="0" style="background-color:#fff;"><tr><td align="center"><table cellpadding="0" cellspacing="0" border="0" style="width:650px"><tr class="layout-full-width" style="background-color:transparent"><![endif]-->
                      <!--[if (mso)|(IE)]><td align="center" width="650" style="background-color:transparent;width:650px; border-top: 0px solid transparent; border-left: 0px solid transparent; border-bottom: 0px solid transparent; border-right: 0px solid transparent;" valign="top"><table width="100%" cellpadding="0" cellspacing="0" border="0"><tr><td style="padding-right: 0px; padding-left: 0px; padding-top:25px; padding-bottom:25px;"><![endif]-->
                      <div
                        class="col num12"
                        style="
                          min-width: 320px;
                          max-width: 650px;
                          display: table-cell;
                          vertical-align: top;
                          width: 650px;
                        "
                      >
                        <div style="width: 100% !important;">
                          <!--[if (!mso)&(!IE)]><!-->
                          <div
                            style="
                              border-top: 0px solid transparent;
                              border-left: 0px solid transparent;
                              border-bottom: 0px solid transparent;
                              border-right: 0px solid transparent;
                              padding-top: 25px;
                              padding-bottom: 25px;
                              padding-right: 0px;
                              padding-left: 0px;
                            "
                          >
                            <!--<![endif]-->
                            <!--[if mso]><table width="100%" cellpadding="0" cellspacing="0" border="0"><tr><td style="padding-right: 25px; padding-left: 25px; padding-top: 25px; padding-bottom: 25px; font-family: Tahoma, Verdana, sans-serif"><![endif]-->
                            <div
                              style="
                                color: #000000;
                                font-family: 'Roboto', Tahoma, Verdana, Segoe,
                                  sans-serif;
                                line-height: 1.5;
                                padding-top: 25px;
                                padding-right: 25px;
                                padding-bottom: 25px;
                                padding-left: 25px;
                              "
                            >
                              <div
                                style="
                                  line-height: 1.5;
                                  font-size: 12px;
                                  font-family: 'Roboto', Tahoma, Verdana, Segoe,
                                    sans-serif;
                                  color: #000000;
                                  mso-line-height-alt: 18px;
                                "
                              >
                                <p
                                  style="
                                    font-size: 16px;
                                    line-height: 1.5;
                                    word-break: break-word;
                                    text-align: left;
                                    font-family: 'Roboto', Tahoma, Verdana, Segoe,
                                      sans-serif;
                                    mso-line-height-alt: 24px;
                                    margin: 0;
                                  "
                                >
                                  <span style="font-size: 16px;"
                                    >{{.HtmlContent}}</span
                                  >
                                </p>
                              </div>
                            </div>
                            <!--[if mso]></td></tr></table><![endif]-->
                            <!--[if (!mso)&(!IE)]><!-->
                          </div>
                          <!--<![endif]-->
                        </div>
                      </div>
                      <!--[if (mso)|(IE)]></td></tr></table><![endif]-->
                      <!--[if (mso)|(IE)]></td></tr></table></td></tr></table><![endif]-->
                    </div>
                  </div>
                </div>
                <div style="background-color: #f2f1f1;">
                  <div
                    class="block-grid"
                    style="
                      margin: 0 auto;
                      min-width: 320px;
                      max-width: 650px;
                      overflow-wrap: break-word;
                      word-wrap: break-word;
                      word-break: break-word;
                      background-color: transparent;
                    "
                  >
                    <div
                      style="
                        border-collapse: collapse;
                        display: table;
                        width: 100%;
                        background-color: transparent;
                      "
                    >
                      <!--[if (mso)|(IE)]><table width="100%" cellpadding="0" cellspacing="0" border="0" style="background-color:#f2f1f1;"><tr><td align="center"><table cellpadding="0" cellspacing="0" border="0" style="width:650px"><tr class="layout-full-width" style="background-color:transparent"><![endif]-->
                      <!--[if (mso)|(IE)]><td align="center" width="650" style="background-color:transparent;width:650px; border-top: 0px solid transparent; border-left: 0px solid transparent; border-bottom: 0px solid transparent; border-right: 0px solid transparent;" valign="top"><table width="100%" cellpadding="0" cellspacing="0" border="0"><tr><td style="padding-right: 0px; padding-left: 0px; padding-top:30px; padding-bottom:30px;"><![endif]-->
                      <div
                        class="col num12"
                        style="
                          min-width: 320px;
                          max-width: 650px;
                          display: table-cell;
                          vertical-align: top;
                          width: 650px;
                        "
                      >
                        <div style="width: 100% !important;">
                          <!--[if (!mso)&(!IE)]><!-->
                          <div
                            style="
                              border-top: 0px solid transparent;
                              border-left: 0px solid transparent;
                              border-bottom: 0px solid transparent;
                              border-right: 0px solid transparent;
                              padding-top: 30px;
                              padding-bottom: 30px;
                              padding-right: 0px;
                              padding-left: 0px;
                            "
                          >
                            <!--<![endif]-->
                            <!--[if mso]><table width="100%" cellpadding="0" cellspacing="0" border="0"><tr><td style="padding-right: 25px; padding-left: 25px; padding-top: 25px; padding-bottom: 25px; font-family: Tahoma, Verdana, sans-serif"><![endif]-->
                            <div
                              style="
                                color: #9a9999;
                                font-family: 'Roboto', Tahoma, Verdana, Segoe,
                                  sans-serif;
                                line-height: 1.5;
                                padding-top: 25px;
                                padding-right: 25px;
                                padding-bottom: 25px;
                                padding-left: 25px;
                              "
                            >
                              <div
                                style="
                                  line-height: 1.5;
                                  font-size: 12px;
                                  font-family: 'Roboto', Tahoma, Verdana, Segoe,
                                    sans-serif;
                                  color: #9a9999;
                                  mso-line-height-alt: 18px;
                                "
                              >
                                <p
                                  style="
                                    font-size: 12px;
                                    line-height: 1.5;
                                    word-break: break-word;
                                    text-align: center;
                                    font-family: 'Roboto', Tahoma, Verdana, Segoe,
                                      sans-serif;
                                    mso-line-height-alt: 18px;
                                    margin: 0;
                                  "
                                >
                                  <span style="font-size: 12px;"
                                    >You will receive emails such as this from time
                                    to time to inform you about account-related
                                    matters.
                                    <br/>
                                    <br/>
                                  </span>
                                </p>
                                <p
                                  style="
                                    font-size: 12px;
                                    line-height: 1.5;
                                    word-break: break-word;
                                    text-align: center;
                                    font-family: 'Roboto', Tahoma, Verdana, Segoe,
                                      sans-serif;
                                    mso-line-height-alt: 18px;
                                    margin: 0;
                                  "
                                >
                                  <span style="font-size: 12px;"
                                    >{{.Company}}, {{.Address}}</span
                                  >
                                </p>
                                <p
                                  style="
                                    font-size: 12px;
                                    line-height: 1.5;
                                    word-break: break-word;
                                    text-align: center;
                                    font-family: 'Roboto', Tahoma, Verdana, Segoe,
                                      sans-serif;
                                    mso-line-height-alt: 18px;
                                    margin: 0;
                                  "
                                >
                                  <span style="font-size: 12px;"
                                    >© Copyright {{.Year}} {{.Company}}</span
                                  >
                                </p>
                              </div>
                            </div>
                            <!--[if mso]></td></tr></table><![endif]-->
                            <!--[if (!mso)&(!IE)]><!-->
                          </div>
                          <!--<![endif]-->
                        </div>
                      </div>
                      <!--[if (mso)|(IE)]></td></tr></table><![endif]-->
                      <!--[if (mso)|(IE)]></td></tr></table></td></tr></table><![endif]-->
                    </div>
                  </div>
                </div>
                <!--[if (mso)|(IE)]></td></tr></table><![endif]-->
              </td>
            </tr>
          </tbody>
        </table>
        <!--[if (IE)]></div><![endif]-->
      </body>
    </html>`

type BrandedEmailData struct {
	Year        string
	Address     string
	Company     string
	HtmlContent string
	LogoLink    string
	Domain      string
	HeaderColor string
}

type BrandedEmailSender struct {
	ProductName   string
	Company       string
	Address       string
	LogoLink      string
	Domain        string
	HeaderColor   string
	EmailTemplate *template.Template
	EmailDialer   *gomail.Dialer
}

func NewBrandedEmailSender(emailDialer *gomail.Dialer) BrandedEmailSender {
	emailTemplate := template.New("Email Template")
	emailTemplate, err := emailTemplate.Parse(rawEmailTemplate)
	if err != nil {
		panic(err)
	}

	return BrandedEmailSender{
		Company:       os.Getenv("COMPANY_NAME"),
		Address:       os.Getenv("COMPANY_ADDRESS"),
		LogoLink:      os.Getenv("LOGO_LINK"),
		Domain:        os.Getenv("APP_DOMAIN"),
		HeaderColor:   os.Getenv("HEADER_COLOR"),
		ProductName:   os.Getenv("APP_NAME"),
		EmailTemplate: emailTemplate,
		EmailDialer:   emailDialer,
	}
}

func (b *BrandedEmailSender) sendPasswordChangedEmail(recipient string) {
	bodyBuffer := new(bytes.Buffer)

	b.EmailTemplate.Execute(bodyBuffer, BrandedEmailData{
		Company:     b.Company,
		Address:     b.Address,
		LogoLink:    b.LogoLink,
		Domain:      b.Domain,
		HeaderColor: b.HeaderColor,
		Year:        strconv.Itoa(time.Now().Year()),
		HtmlContent: "You just changed the password of your " + b.ProductName +
			" account. If this was a mistake contact us to avoid losing access to your account.<br/><br/>Cheers,<br/>The " +
			b.ProductName + " team",
	})

	m := gomail.NewMessage()
	m.SetHeader("From", "binder@baida.dev")
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", "Password successfully changed")
	m.SetBody("text/html", bodyBuffer.String())
	m.AddAlternative("text/plain", "You just changed the password of your "+b.ProductName+
		" account. If this was a mistake contact us to avoid losing access to your account.\n\nCheers,\nThe "+
		b.ProductName+" team")

	if err := b.EmailDialer.DialAndSend(m); err != nil {
		panic(err)
	}
}
