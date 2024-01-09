package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-mail/mail"
	"github.com/joho/godotenv"
)

func SendEmailForContest(toEmail string, contest_id string, account string, user_password string, user_investor_pass string, promo_code string) error {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file")
	}

	username := os.Getenv("EMAIL_ADDRESS")
	password := os.Getenv("EMAIL_PASSWORD")
	trading_server := os.Getenv("EMAIL_SERVER_TRADING")
	trading_platform := os.Getenv("EMAIL_PLATFORM_TRADING")
	// Gmail SMTP server and port
	smtpHost := "smtp.zoho.com"
	smtpPort := 587

	// Email subject and body
	subject := fmt.Sprintf("Trading account details for the competition: %s", contest_id)
	body := `<body>
						<div id=":mw" class="ii gt">
							<div id=":mv" class="a3s aiL msg8455739578090132613"><u></u>
								<div style="margin:0!important;padding:0!important;background-color:#e6eaed" bgcolor="#E6EAED">
									<div style="display:none;font-size:1px;color:#8957FF;line-height:1px;font-family:Open Sans,Helvetica,Arial,sans-serif;max-height:0px;max-width:0px;opacity:0;overflow:hidden">
										Save your ` + trading_platform + ` account details <br> and start trading.
									</div>
									<table border="0" cellpadding="0" cellspacing="0" width="100%">
										<tbody>
											<tr>
												<td align="center" style="background-color:#e6eaed;padding:30px!important" bgcolor="#E6EAED">
													<table align="center" border="0" cellpadding="0" cellspacing="0" width="600" style="min-width:600px!important">
														<tbody>
															<tr>
																<td align="center" style="padding:42px 40px 0 40px;border-radius:3px 3px 0px 0px" bgcolor="#FFFFFF">
																	<table width="520" cellpadding="0" cellspacing="0" border="0" align="center">
																		<tbody>
																			<tr>
																				<td align="center">
																					<table cellpadding="0" cellspacing="0" border="0">
																						<tbody>
																							<tr>
																								<td align="center">
																									<a href="https://fxchampionship.com" target="_blank" data-saferedirecturl="https://fxchampionship.com"><img
																											src="https://crm.fxchampionship.com/src/assets/images/logos/android-chrome-512x512.png"
																											width="160" height="160"
																											style="margin:0;padding:0;border:none;display:block"
																											border="0"
																											alt="FXChampionship"
																											class="CToWUd"
																											data-bit="iit">
																									</a>
																								</td>
																							</tr>
																						</tbody>
																					</table>
																				</td>
																			</tr>
																		</tbody>
																	</table>
																</td>
															</tr>
															<tr>
																<td align="center" style="padding:0 40px 0 40px;background-color:#ffffff" bgcolor="#FFFFFF">
																	<table align="center" border="0" cellpadding="0" cellspacing="0" width="100%">
																		<tbody>
																			<tr>
																				<td align="centar" style="font:300 44px/48px roboto,helvetica,sans-serif;color:#000000;padding:40px 0 60px;margin:0;text-align:center">
																					Save your ` + trading_platform + ` account details and start trading.
																				</td>
																			</tr>
																			<tr>
																				<td>
																					<div style="background:#f0f5f8;border:1px solid #0a95ff;border-radius:4px;padding:40px">
																						<table align="center" border="0" cellpadding="0" cellspacing="0" width="500" bgcolor="#f0f5f8">
																							<tbody>
																								<tr>
																									<td
																										style="font:300 20px/28px roboto,helvetica,sans-serif;color:#2d3436;padding:0 10px 5px 0;margin:0">
																										Login:</td>
																									<td
																										style="font:400 20px/28px roboto,helvetica,sans-serif;color:#000000;padding:0 0 5px;margin:0">
																										` + account + `</td>
																								</tr>
																								<tr>
																									<td
																										style="font:300 20px/28px roboto,helvetica,sans-serif;color:#2d3436;padding:0 10px 5px 0;margin:0">
																										Password:</td>
																									<td
																										style="font:400 20px/28px roboto,helvetica,sans-serif;color:#000000;padding:0 0 5px;margin:0">
																										<u></u>` + user_password + `<u></u>
																									</td>
																								</tr>
																								<tr>
																									<td
																										style="font:300 20px/28px roboto,helvetica,sans-serif;color:#2d3436;padding:0 10px 5px 0;margin:0">
																										Investor Password:</td>
																									<td
																										style="font:400 20px/28px roboto,helvetica,sans-serif;color:#000000;padding:0 0 5px;margin:0">
																										<u></u>` + user_investor_pass + `<u></u>
																									</td>
																								</tr>
																								<tr>
																									<td
																										style="font:300 20px/28px roboto,helvetica,sans-serif;color:#2d3436;padding:0 10px 0 0;margin:0">
																										Server:</td>
																									<td
																										style="font:400 20px/28px roboto,helvetica,sans-serif;color:#000000;padding:0;margin:0">
																										` + trading_server + `</td>
																								</tr>`
	if promo_code != "" {
		body += `<tr>
																									<td
																										style="font:300 20px/28px roboto,helvetica,sans-serif;color:#2d3436;padding:0 10px 5px 0;margin:0">
																										Promotion code:</td>
																									<td
																										style="font:400 20px/28px roboto,helvetica,sans-serif;color:#000000;padding:0 0 5px;margin:0">
																										<u></u>` + promo_code + `<u></u>
																									</td>
																								</tr>`
	}
	body += `</tbody>
																						</table>
																					</div>
																				</td>
																			</tr>
																			<tr>
																				<td align="left" style="font:300 20px/28px roboto,helvetica,sans-serif;color:#636e72;padding:40px 0 0;margin:0">Use this information to access the trading platform you selected during the registration process.
																				</td>
																			</tr>
																			<tr>
																				<td align="left" style="font:300 20px/28px roboto,helvetica,sans-serif;color:#636e72;padding:25px 0 0;margin:0">
																					If you don't have it yet, please proceed with the
																					installation here
																				</td>
																			</tr>
																		</tbody>
																	</table>
																</td>
															</tr>
															<tr>
																<td align="center" valign="middle" width="100%" bgcolor="#FFFFFF" style="padding:40px 0 0px;text-align:center">

																	<div>
																		<a href="https://www.cptmarkets.com/platform/metatrader4"
																			style="width:412px;background-color:#00af52;background-repeat:no-repeat;border:1px solid #00af52;border-radius:2px;color:#ffffff;display:inline-block;font:700 20px/26px roboto,helvetica,sans-serif;text-decoration:none;padding:20px 0;text-align:center"
																			target="_blank"
																			data-saferedirecturl="https://www.cptmarkets.com/platform/metatrader4">
																			Download your platform <br>
																			<span
																				style="color:#9cd6b7;font-size:14px">Web/Mac/Win/Android/iOS</span>
																		</a>
																	</div>
																</td>
															</tr>
															<tr>
																<td align="center" style="padding:0 40px 0 40px;background-color:#ffffff"
																	bgcolor="#FFFFFF">
																	<table align="center" border="0" cellpadding="0" cellspacing="0"
																		width="100%">
																		<tbody>
																			<tr>
																				<td align="center" style="font:300 44px/50px roboto,Arial,helvetica,sans-serif;color:#000000;padding:30px 0 60px;margin:0;text-align:center">
																					What should I do next?
																				</td>
																			</tr>
																			<tr>
																				<td align="center" style="font:300 20px/30px roboto,Arial,helvetica,sans-serif;color:#000000;margin:0;text-align:center">
																					Log in to the ` + trading_platform + ` platform, start your first
																					trade, and show everyone that you're a real formidable
																					competitor!
																				</td>
																			</tr>
																		</tbody>
																	</table>
																</td>
															</tr>
															<tr>
																<td align="center" style="padding:40px 40px 40px;background-color:#ffffff"
																	bgcolor="#FFFFFF">
																	<table align="center" border="0" cellpadding="0" cellspacing="0"
																		width="100%">
																		<tbody>
																			<tr>
																				<td align="left" style="color:#636e72;font:400 16px/22px roboto,helvetica,sans-serif;padding:40px 40px 25px;border-top:1px solid #e7e9ec">
																					If you have any questions, please don't hesitate to
																					contact our 24/7 customer support department at: <a
																						href="mailto:support@fxchampionship.com"
																						style="color:#0a95ff;font-weight:500;direction:ltr;text-decoration:none"
																						target="_blank">support@<span
																							class="il">fxchampionship</span>.com</a><br> or
																					by phone at<a href="tel:%20+44(0)2031515550"
																						style="color:#000000;font-weight:500;direction:ltr;text-decoration:none"
																						target="_blank"> +84 919 720 567</a>
																				</td>
																			</tr>
																		</tbody>
																	</table>
																</td>
															</tr>
														</tbody>
													</table>
												</td>
											</tr>
										</tbody>
									</table>
									</div>
								</div>
							</div>
						</div>
					</body>`

	m := mail.NewMessage()
	m.SetHeader("From", username)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := mail.NewDialer(smtpHost, smtpPort, username, password)
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		return err
	} else {
		fmt.Printf("Send email to: %s\n", toEmail)
	}
	return nil
}

func SendEmailForRegister(toEmail string, account string, user_password string) error {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file")
	}

	username := os.Getenv("EMAIL_ADDRESS")
	password := os.Getenv("EMAIL_PASSWORD")
	// Gmail SMTP server and port
	smtpHost := "smtp.zoho.com"
	smtpPort := 587

	// Email subject and body
	subject := "Welcome to FXChampionship - A Gathering Place for World-Class Traders!"
	body := `<body>
    <div id=":mw" class="ii gt">
        <div id=":mv" class="a3s aiL msg8455739578090132613"><u></u>
            <div style="margin:0!important;padding:0!important;background-color:#e6eaed" bgcolor="#E6EAED">
                <div
                    style="display:none;font-size:1px;color:#8957FF;line-height:1px;font-family:Open Sans,Helvetica,Arial,sans-serif;max-height:0px;max-width:0px;opacity:0;overflow:hidden">
                    Welcome to FXChampionship - A Gathering Place for World-Class Traders!
                </div>
                <table border="0" cellpadding="0" cellspacing="0" width="100%">
                    <tbody>
                        <tr>
                            <td align="center" style="background-color:#e6eaed;padding:30px!important"
                                bgcolor="#E6EAED">
                                <table align="center" border="0" cellpadding="0" cellspacing="0" width="600"
                                    style="min-width:600px!important">
                                    <tbody>
                                        <tr>
                                            <td align="center"
                                                style="padding:42px 40px 0 40px;border-radius:3px 3px 0px 0px"
                                                bgcolor="#FFFFFF">
                                                <table width="520" cellpadding="0" cellspacing="0" border="0"
                                                    align="center">
                                                    <tbody>
                                                        <tr>
                                                            <td align="center">
                                                                <table cellpadding="0" cellspacing="0" border="0">
                                                                    <tbody>
                                                                        <tr>
                                                                            <td align="center">
                                                                                <a href="https://fxchampionship.com"
                                                                                    target="_blank"
                                                                                    data-saferedirecturl="https://fxchampionship.com"><img
                                                                                        src="https://crm.fxchampionship.com/src/assets/images/logos/android-chrome-512x512.png"
                                                                                        width="160" height="160"
                                                                                        style="margin:0;padding:0;border:none;display:block"
                                                                                        border="0"
                                                                                        alt="FxChampionship"
                                                                                        class="CToWUd"
                                                                                        data-bit="iit"></a>
                                                                            </td>
                                                                        </tr>
                                                                    </tbody>
                                                                </table>
                                                            </td>
                                                        </tr>
                                                    </tbody>
                                                </table>
                                            </td>
                                        </tr>
                                        <tr>
                                            <td align="center" style="padding:0 40px 0 40px;background-color:#ffffff"
                                                bgcolor="#FFFFFF">
                                                <table align="center" border="0" cellpadding="0" cellspacing="0"
                                                    width="100%">
                                                    <tbody>
                                                        <tr>
                                                            <td align="centar"
                                                                style="font:300 40px/44px roboto,helvetica,sans-serif;color:#000000;padding-top: 10px;margin:0;text-align:center">
                                                                Welcome to FXChampionship
                                                            </td>
                                                        </tr>
                                                        <tr>
                                                            <td align="centar"
                                                                style="font:300 22px/28px roboto,helvetica,sans-serif;color:#000000;padding-bottom: 50px;padding-top: 10px;margin:0;text-align:center">
                                                                A Gathering Place for
                                                                World-Class Traders!
                                                            </td>
                                                        </tr>
                                                        <tr>
                                                            <td>
                                                                <div
                                                                    style="background:#f0f5f8;border:1px solid #0a95ff;border-radius:4px;padding:40px">
                                                                    <table align="center" border="0" cellpadding="0"
                                                                        cellspacing="0" width="500" bgcolor="#f0f5f8">
                                                                        <tbody>
                                                                            <tr>
                                                                                <td
                                                                                    style="font:300 20px/28px roboto,helvetica,sans-serif;color:#2d3436;padding:0 10px 5px 0;margin:0">
                                                                                    Email:</td>
                                                                                <td
                                                                                    style="font:400 20px/28px roboto,helvetica,sans-serif;color:#000000;padding:0 0 5px;margin:0">
                                                                                    ` + account + `</td>
                                                                            </tr>
                                                                            <tr>
                                                                                <td
                                                                                    style="font:300 20px/28px roboto,helvetica,sans-serif;color:#2d3436;padding:0 10px 5px 0;margin:0">
                                                                                    Password:</td>
                                                                                <td
                                                                                    style="font:400 20px/28px roboto,helvetica,sans-serif;color:#000000;padding:0 0 5px;margin:0">
                                                                                    <u></u>` + user_password + `<u></u>
                                                                                </td>
                                                                            </tr>
                                                                        </tbody>
                                                                    </table>
                                                                </div>
                                                            </td>
                                                        </tr>
                                                        <tr>
                                                            <td align="left"
                                                                style="font:300 20px/28px roboto,helvetica,sans-serif;color:#636e72;padding:40px 0 0;margin:0">
                                                                Use this information to access the trading platform you
                                                                selected during the registration process.
                                                            </td>
                                                        </tr>
                                                        <tr>
                                                            <td align="left"
                                                                style="font:300 20px/28px roboto,helvetica,sans-serif;color:#636e72;padding:25px 0 0;margin:0">
                                                                If you don't have it yet, please proceed with the
                                                                installation here
                                                            </td>
                                                        </tr>
                                                    </tbody>
                                                </table>
                                            </td>
                                        </tr>
                                        <tr>
                                            <td align="center" valign="middle" width="100%" bgcolor="#FFFFFF"
                                                style="padding:40px 0 0px;text-align:center">

                                                <div>
                                                    <a href="https://www.cptmarkets.com/platform/metatrader4"
                                                        style="width:412px;background-color:#00af52;background-repeat:no-repeat;border:1px solid #00af52;border-radius:2px;color:#ffffff;display:inline-block;font:700 20px/26px roboto,helvetica,sans-serif;text-decoration:none;padding:20px 0;text-align:center"
                                                        target="_blank"
                                                        data-saferedirecturl="https://www.cptmarkets.com/platform/metatrader4">
                                                        Download your platform <br>
                                                        <span
                                                            style="color:#9cd6b7;font-size:14px">Web/Mac/Win/Android/iOS</span>
                                                    </a>
                                                </div>
                                            </td>
                                        </tr>
                                        <tr>
                                            <td align="center" style="padding:0 40px 0 40px;background-color:#ffffff"
                                                bgcolor="#FFFFFF">
                                                <table align="center" border="0" cellpadding="0" cellspacing="0"
                                                    width="100%">
                                                    <tbody>
                                                        <tr>
                                                            <td align="center"
                                                                style="font:300 44px/50px roboto,Arial,helvetica,sans-serif;color:#000000;padding:30px 0 20px;margin:0;text-align:center">
                                                                What should I do next?
                                                            </td>
                                                        </tr>
                                                        <tr>
                                                            <td align="center"
                                                                style="font:300 20px/30px roboto,Arial,helvetica,sans-serif;color:#000000;margin:0;text-align:center">
                                                                Join a competition and emerge victorious!
                                                            </td>
                                                        </tr>
                                                    </tbody>
                                                </table>
                                            </td>
                                        </tr>

                                        <tr>
                                            <td align="center" style="padding:40px 40px 40px;background-color:#ffffff"
                                                bgcolor="#FFFFFF">
                                                <table align="center" border="0" cellpadding="0" cellspacing="0"
                                                    width="100%">
                                                    <tbody>
                                                        <tr>
                                                            <td align="left"
                                                                style="color:#636e72;font:400 16px/22px roboto,helvetica,sans-serif;padding:40px 40px 25px;border-top:1px solid #e7e9ec">
                                                                If you have any questions, please don't hesitate to
                                                                contact our 24/7 customer support department at: <a
                                                                    href="mailto:support@fxchampionship.com"
                                                                    style="color:#0a95ff;font-weight:500;direction:ltr;text-decoration:none"
                                                                    target="_blank">support@<span
                                                                        class="il">fxchampionship</span>.com</a><br> or
                                                                by phone at<a href="tel:%20+44(0)2031515550"
                                                                    style="color:#000000;font-weight:500;direction:ltr;text-decoration:none"
                                                                    target="_blank"> +84 919 720 567</a>
                                                            </td>
                                                        </tr>
                                                    </tbody>
                                                </table>
                                            </td>
                                        </tr>
                                    </tbody>
                                </table>
                            </td>
                        </tr>
                    </tbody>
                </table>
                <div class="yj6qo"></div>
                <div class="adL">
                </div>
                <div style="display:none;white-space:nowrap;font:15px courier;line-height:0" class="adL"> &nbsp; &nbsp;
                    &nbsp; &nbsp;
                    &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;
                    &nbsp;
                    &nbsp;
                    &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;</div>
                <div class="adL">


                </div>
            </div>
            <div class="adL">


            </div>
        </div>
    </div>
</body>`

	m := mail.NewMessage()
	m.SetHeader("From", username)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	// m.Attach("/Users/vngoldstreet/Downloads/thantai-goldenfund.jpeg")

	d := mail.NewDialer(smtpHost, smtpPort, username, password)
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		return err
	} else {
		fmt.Printf("Send email to: %s\n", toEmail)
	}
	return nil
}

func SendEmailForResetPassword(toEmail string, account string, user_password string) error {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file")
	}

	username := os.Getenv("EMAIL_ADDRESS")
	password := os.Getenv("EMAIL_PASSWORD")
	// Gmail SMTP server and port
	smtpHost := "smtp.zoho.com"
	smtpPort := 587

	// Email subject and body
	subject := "FXChampionship - Reset Password"
	body := `<body>
    <div id=":mw" class="ii gt">
        <div id=":mv" class="a3s aiL msg8455739578090132613"><u></u>
            <div style="margin:0!important;padding:0!important;background-color:#e6eaed" bgcolor="#E6EAED">
                <div
                    style="display:none;font-size:1px;color:#8957FF;line-height:1px;font-family:Open Sans,Helvetica,Arial,sans-serif;max-height:0px;max-width:0px;opacity:0;overflow:hidden">
                    FXChampionship - Reset Password
                </div>
                <table border="0" cellpadding="0" cellspacing="0" width="100%">
                    <tbody>
                        <tr>
                            <td align="center" style="background-color:#e6eaed;padding:30px!important"
                                bgcolor="#E6EAED">
                                <table align="center" border="0" cellpadding="0" cellspacing="0" width="600"
                                    style="min-width:600px!important">
                                    <tbody>
                                        <tr>
                                            <td align="center"
                                                style="padding:42px 40px 0 40px;border-radius:3px 3px 0px 0px"
                                                bgcolor="#FFFFFF">
                                                <table width="520" cellpadding="0" cellspacing="0" border="0"
                                                    align="center">
                                                    <tbody>
                                                        <tr>
                                                            <td align="center">
                                                                <table cellpadding="0" cellspacing="0" border="0">
                                                                    <tbody>
                                                                        <tr>
                                                                            <td align="center">
                                                                                <a href="https://fxchampionship.com"
                                                                                    target="_blank"
                                                                                    data-saferedirecturl="https://fxchampionship.com"><img
                                                                                        src="https://crm.fxchampionship.com/src/assets/images/logos/android-chrome-512x512.png"
                                                                                        width="160" height="160"
                                                                                        style="margin:0;padding:0;border:none;display:block"
                                                                                        border="0"
                                                                                        alt="FxChampionship"
                                                                                        class="CToWUd"
                                                                                        data-bit="iit"></a>
                                                                            </td>
                                                                        </tr>
                                                                    </tbody>
                                                                </table>
                                                            </td>
                                                        </tr>
                                                    </tbody>
                                                </table>
                                            </td>
                                        </tr>
                                        <tr>
                                            <td align="center" style="padding:0 40px 0 40px;background-color:#ffffff"
                                                bgcolor="#FFFFFF">
                                                <table align="center" border="0" cellpadding="0" cellspacing="0"
                                                    width="100%">
                                                    <tbody>
                                                        <tr>
                                                            <td align="centar"
                                                                style="font:300 40px/44px roboto,helvetica,sans-serif;color:#000000;padding-top: 10px;margin:0;text-align:center">
                                                                Welcome to FXChampionship
                                                            </td>
                                                        </tr>
                                                        <tr>
                                                            <td align="centar"
                                                                style="font:300 22px/28px roboto,helvetica,sans-serif;color:#000000;padding-bottom: 50px;padding-top: 10px;margin:0;text-align:center">
                                                                A Gathering Place for
                                                                World-Class Traders!
                                                            </td>
                                                        </tr>
                                                        <tr>
                                                            <td>
                                                                <div
                                                                    style="background:#f0f5f8;border:1px solid #0a95ff;border-radius:4px;padding:40px">
                                                                    <table align="center" border="0" cellpadding="0"
                                                                        cellspacing="0" width="500" bgcolor="#f0f5f8">
                                                                        <tbody>
                                                                            <tr>
                                                                                <td
                                                                                    style="font:300 20px/28px roboto,helvetica,sans-serif;color:#2d3436;padding:0 10px 5px 0;margin:0">
                                                                                    Email:</td>
                                                                                <td
                                                                                    style="font:400 20px/28px roboto,helvetica,sans-serif;color:#000000;padding:0 0 5px;margin:0">
                                                                                    ` + account + `</td>
                                                                            </tr>
                                                                            <tr>
                                                                                <td
                                                                                    style="font:300 20px/28px roboto,helvetica,sans-serif;color:#2d3436;padding:0 10px 5px 0;margin:0">
                                                                                    New Password:</td>
                                                                                <td
                                                                                    style="font:400 20px/28px roboto,helvetica,sans-serif;color:#000000;padding:0 0 5px;margin:0">
                                                                                    <u></u>` + user_password + `<u></u>
                                                                                </td>
                                                                            </tr>
                                                                        </tbody>
                                                                    </table>
                                                                </div>
                                                            </td>
                                                        </tr>
                                                        <tr>
                                                            <td align="left"
                                                                style="font:300 20px/28px roboto,helvetica,sans-serif;color:#636e72;padding:40px 0 0;margin:0">
                                                                Use this information to access the trading platform you
                                                                selected during the registration process.
                                                            </td>
                                                        </tr>
                                                        <tr>
                                                            <td align="left"
                                                                style="font:300 20px/28px roboto,helvetica,sans-serif;color:#636e72;padding:25px 0 0;margin:0">
                                                                If you don't have it yet, please proceed with the
                                                                installation here
                                                            </td>
                                                        </tr>
                                                    </tbody>
                                                </table>
                                            </td>
                                        </tr>
                                        <tr>
                                            <td align="center" valign="middle" width="100%" bgcolor="#FFFFFF"
                                                style="padding:40px 0 0px;text-align:center">

                                                <div>
                                                    <a href="https://www.cptmarkets.com/platform/metatrader4"
                                                        style="width:412px;background-color:#00af52;background-repeat:no-repeat;border:1px solid #00af52;border-radius:2px;color:#ffffff;display:inline-block;font:700 20px/26px roboto,helvetica,sans-serif;text-decoration:none;padding:20px 0;text-align:center"
                                                        target="_blank"
                                                        data-saferedirecturl="https://www.cptmarkets.com/platform/metatrader4">
                                                        Download your platform <br>
                                                        <span
                                                            style="color:#9cd6b7;font-size:14px">Web/Mac/Win/Android/iOS</span>
                                                    </a>
                                                </div>
                                            </td>
                                        </tr>
                                        <tr>
                                            <td align="center" style="padding:0 40px 0 40px;background-color:#ffffff"
                                                bgcolor="#FFFFFF">
                                                <table align="center" border="0" cellpadding="0" cellspacing="0"
                                                    width="100%">
                                                    <tbody>
                                                        <tr>
                                                            <td align="center"
                                                                style="font:300 44px/50px roboto,Arial,helvetica,sans-serif;color:#000000;padding:30px 0 20px;margin:0;text-align:center">
                                                                What should I do next?
                                                            </td>
                                                        </tr>
                                                        <tr>
                                                            <td align="center"
                                                                style="font:300 20px/30px roboto,Arial,helvetica,sans-serif;color:#000000;margin:0;text-align:center">
                                                                Join a competition and emerge victorious!
                                                            </td>
                                                        </tr>
                                                    </tbody>
                                                </table>
                                            </td>
                                        </tr>

                                        <tr>
                                            <td align="center" style="padding:40px 40px 40px;background-color:#ffffff"
                                                bgcolor="#FFFFFF">
                                                <table align="center" border="0" cellpadding="0" cellspacing="0"
                                                    width="100%">
                                                    <tbody>
                                                        <tr>
                                                            <td align="left"
                                                                style="color:#636e72;font:400 16px/22px roboto,helvetica,sans-serif;padding:40px 40px 25px;border-top:1px solid #e7e9ec">
                                                                If you have any questions, please don't hesitate to
                                                                contact our 24/7 customer support department at: <a
                                                                    href="mailto:support@fxchampionship.com"
                                                                    style="color:#0a95ff;font-weight:500;direction:ltr;text-decoration:none"
                                                                    target="_blank">support@<span
                                                                        class="il">fxchampionship</span>.com</a><br> or
                                                                by phone at<a href="tel:%20+44(0)2031515550"
                                                                    style="color:#000000;font-weight:500;direction:ltr;text-decoration:none"
                                                                    target="_blank"> +84 919 720 567</a>
                                                            </td>
                                                        </tr>
                                                    </tbody>
                                                </table>
                                            </td>
                                        </tr>
                                    </tbody>
                                </table>
                            </td>
                        </tr>
                    </tbody>
                </table>
                <div class="yj6qo"></div>
                <div class="adL">
                </div>
                <div style="display:none;white-space:nowrap;font:15px courier;line-height:0" class="adL"> &nbsp; &nbsp;
                    &nbsp; &nbsp;
                    &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;
                    &nbsp;
                    &nbsp;
                    &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;</div>
                <div class="adL">


                </div>
            </div>
            <div class="adL">


            </div>
        </div>
    </div>
</body>`

	m := mail.NewMessage()
	m.SetHeader("From", username)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	// m.Attach("/Users/vngoldstreet/Downloads/thantai-goldenfund.jpeg")

	d := mail.NewDialer(smtpHost, smtpPort, username, password)
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		return err
	} else {
		fmt.Printf("Send email to: %s\n", toEmail)
	}
	return nil
}
