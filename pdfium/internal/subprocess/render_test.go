package subprocess_test

import (
	"bytes"
	"encoding/gob"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"os"

	"github.com/klippa-app/go-pdfium/pdfium/internal/subprocess"
	"github.com/klippa-app/go-pdfium/pdfium/pdfium_errors"
	"github.com/klippa-app/go-pdfium/pdfium/requests"
	"github.com/klippa-app/go-pdfium/pdfium/responses"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Render", func() {
	pdfium := subprocess.Pdfium{}

	Context("no document", func() {
		When("is opened", func() {
			Context("GetPageSize()", func() {
				It("returns an error", func() {
					pageSize, err := pdfium.GetPageSize(&requests.GetPageSize{
						Page: 0,
					})
					Expect(err).To(MatchError("no current document"))
					Expect(pageSize).To(BeNil())
				})
			})

			Context("GetPageSizeInPixels()", func() {
				It("returns an error", func() {
					pageSize, err := pdfium.GetPageSizeInPixels(&requests.GetPageSizeInPixels{
						Page: 0,
						DPI:  100,
					})
					Expect(err).To(MatchError("no current document"))
					Expect(pageSize).To(BeNil())
				})
			})

			Context("RenderPageInDPI()", func() {
				It("returns an error", func() {
					renderedPage, err := pdfium.RenderPageInDPI(&requests.RenderPageInDPI{
						Page: 0,
						DPI:  300,
					})
					Expect(err).To(MatchError("no current document"))
					Expect(renderedPage).To(BeNil())
				})
			})

			Context("RenderPageInDPI()", func() {
				It("returns an error", func() {
					renderedPage, err := pdfium.RenderPagesInDPI(&requests.RenderPagesInDPI{
						Pages: []requests.RenderPageInDPI{
							{
								Page: 0,
								DPI:  300,
							},
						},
						Padding: 50,
					})
					Expect(err).To(MatchError("no current document"))
					Expect(renderedPage).To(BeNil())
				})
			})

			Context("RenderPageInPixels()", func() {
				It("returns an error", func() {
					renderedPage, err := pdfium.RenderPageInPixels(&requests.RenderPageInPixels{
						Page:   0,
						Width:  2000,
						Height: 2000,
					})
					Expect(err).To(MatchError("no current document"))
					Expect(renderedPage).To(BeNil())
				})
			})

			Context("RenderPagesInPixels()", func() {
				It("returns an error", func() {
					renderedPage, err := pdfium.RenderPagesInPixels(&requests.RenderPagesInPixels{
						Pages: []requests.RenderPageInPixels{
							{
								Page:   0,
								Width:  2000,
								Height: 2000,
							},
						},
						Padding: 50,
					})
					Expect(err).To(MatchError("no current document"))
					Expect(renderedPage).To(BeNil())
				})
			})
		})
	})

	Context("a normal PDF file", func() {
		BeforeEach(func() {
			pdfData, _ := ioutil.ReadFile("./testdata/test.pdf")
			pdfium.OpenDocument(&requests.OpenDocument{
				File: &pdfData,
			})
		})

		AfterEach(func() {
			pdfium.Close()
		})

		When("is opened", func() {
			Context("when an invalid page is given", func() {
				Context("GetPageSize()", func() {
					It("returns an error", func() {
						pageSize, err := pdfium.GetPageSize(&requests.GetPageSize{
							Page: 1,
						})
						Expect(err).To(MatchError(pdfium_errors.ErrPage))
						Expect(pageSize).To(BeNil())
					})
				})

				Context("GetPageSizeInPixels()", func() {
					It("returns an error", func() {
						pageSize, err := pdfium.GetPageSizeInPixels(&requests.GetPageSizeInPixels{
							Page: 1,
							DPI:  100,
						})
						Expect(err).To(MatchError(pdfium_errors.ErrPage))
						Expect(pageSize).To(BeNil())
					})
				})

				Context("RenderPageInDPI()", func() {
					It("returns an error", func() {
						renderedPage, err := pdfium.RenderPageInDPI(&requests.RenderPageInDPI{
							Page: 1,
							DPI:  300,
						})
						Expect(err).To(MatchError(pdfium_errors.ErrPage))
						Expect(renderedPage).To(BeNil())
					})
				})

				Context("RenderPagesInDPI()", func() {
					It("returns an error", func() {
						renderedPage, err := pdfium.RenderPagesInDPI(&requests.RenderPagesInDPI{
							Pages: []requests.RenderPageInDPI{
								{
									Page: 1,
									DPI:  300,
								},
							},
							Padding: 50,
						})
						Expect(err).To(MatchError(pdfium_errors.ErrPage))
						Expect(renderedPage).To(BeNil())
					})
				})

				Context("RenderPageInPixels()", func() {
					It("returns an error", func() {
						renderedPage, err := pdfium.RenderPageInPixels(&requests.RenderPageInPixels{
							Page:   1,
							Width:  2000,
							Height: 2000,
						})
						Expect(err).To(MatchError(pdfium_errors.ErrPage))
						Expect(renderedPage).To(BeNil())
					})
				})
			})

			Context("RenderPagesInPixels()", func() {
				It("returns an error", func() {
					renderedPage, err := pdfium.RenderPagesInPixels(&requests.RenderPagesInPixels{
						Pages: []requests.RenderPageInPixels{
							{
								Page:   1,
								Width:  2000,
								Height: 2000,
							},
						},
						Padding: 50,
					})
					Expect(err).To(MatchError(pdfium_errors.ErrPage))
					Expect(renderedPage).To(BeNil())
				})
			})

			Context("when the page size is requested", func() {
				Context("in points", func() {
					It("returns the correct amount of points", func() {
						pageSize, err := pdfium.GetPageSize(&requests.GetPageSize{
							Page: 0,
						})
						Expect(err).To(BeNil())
						Expect(pageSize).To(Equal(&responses.GetPageSize{
							Width:  595.2755737304688,
							Height: 841.8897094726562,
						}))
					})
				})

				Context("in pixels", func() {
					Context("with no DPI", func() {
						It("returns an error", func() {
							pageSize, err := pdfium.GetPageSizeInPixels(&requests.GetPageSizeInPixels{
								Page: 0,
							})
							Expect(err).To(MatchError("no DPI given"))
							Expect(pageSize).To(BeNil())
						})
					})

					Context("width DPI 100", func() {
						It("returns the right amount of pixels and point to pixel ratio", func() {
							pageSize, err := pdfium.GetPageSizeInPixels(&requests.GetPageSizeInPixels{
								Page: 0,
								DPI:  100,
							})
							Expect(err).To(BeNil())
							Expect(pageSize).To(Equal(&responses.GetPageSizeInPixels{
								Width:             827,
								Height:            1170,
								PointToPixelRatio: 1.3888888888888888,
							}))
						})
					})

					Context("width DPI 300", func() {
						It("returns the right amount of pixels and point to pixel ratio", func() {
							pageSize, err := pdfium.GetPageSizeInPixels(&requests.GetPageSizeInPixels{
								Page: 0,
								DPI:  300,
							})
							Expect(err).To(BeNil())
							Expect(pageSize).To(Equal(&responses.GetPageSizeInPixels{
								Width:             2481,
								Height:            3508,
								PointToPixelRatio: 4.166666666666667,
							}))
						})
					})
				})
			})

			Context("the page is rendered", func() {
				Context("in points", func() {
					Context("with no DPI", func() {
						It("returns an error", func() {
							renderedPage, err := pdfium.RenderPageInDPI(&requests.RenderPageInDPI{
								Page: 0,
							})
							Expect(err).To(MatchError("no DPI given"))
							Expect(renderedPage).To(BeNil())
						})
					})

					Context("width DPI 100", func() {
						It("returns the right image, point to pixel ratio and resolution", func() {
							renderedPage, err := pdfium.RenderPageInDPI(&requests.RenderPageInDPI{
								Page: 0,
								DPI:  100,
							})
							Expect(err).To(BeNil())
							Expect(renderedPage).To(Equal(&responses.RenderPage{
								Image:             loadPrerenderedImage("./testdata/render_testpdf_dpi_100.gob", renderedPage.Image),
								PointToPixelRatio: 1.3888888888888888,
							}))
							Expect(renderedPage.Image.Bounds().Size().X).To(Equal(827))
							Expect(renderedPage.Image.Bounds().Size().Y).To(Equal(1170))
						})
					})

					Context("width DPI 300", func() {
						It("returns the right image, point to pixel ratio and resolution", func() {
							renderedPage, err := pdfium.RenderPageInDPI(&requests.RenderPageInDPI{
								Page: 0,
								DPI:  300,
							})
							Expect(err).To(BeNil())
							Expect(renderedPage).To(Equal(&responses.RenderPage{
								Image:             loadPrerenderedImage("./testdata/render_testpdf_dpi_300.gob", renderedPage.Image),
								PointToPixelRatio: 4.166666666666667,
							}))

							Expect(renderedPage.Image.Bounds().Size().X).To(Equal(2481))
							Expect(renderedPage.Image.Bounds().Size().Y).To(Equal(3508))
						})
					})
				})

				Context("in pixels", func() {
					Context("with no width or height given", func() {
						It("returns an error", func() {
							renderedPage, err := pdfium.RenderPageInPixels(&requests.RenderPageInPixels{
								Page: 0,
							})
							Expect(err).To(MatchError("no width or height given"))
							Expect(renderedPage).To(BeNil())
						})
					})

					Context("with only the width given", func() {
						It("returns the right image, point to pixel ratio and resolution", func() {
							renderedPage, err := pdfium.RenderPageInPixels(&requests.RenderPageInPixels{
								Page:  0,
								Width: 2000,
							})

							Expect(err).To(BeNil())
							Expect(renderedPage).To(Equal(&responses.RenderPage{
								Image:             loadPrerenderedImage("./testdata/render_testpdf_pixels_2000x0.gob", renderedPage.Image),
								PointToPixelRatio: 3.3597884547259587,
							}))
							Expect(renderedPage.Image.Bounds().Size().X).To(Equal(2000))
							Expect(renderedPage.Image.Bounds().Size().Y).To(Equal(2829))
						})
					})

					Context("with only the height given", func() {
						It("returns the right image, point to pixel ratio and resolution", func() {
							renderedPage, err := pdfium.RenderPageInPixels(&requests.RenderPageInPixels{
								Page:   0,
								Height: 2000,
							})

							Expect(err).To(BeNil())
							Expect(renderedPage).To(Equal(&responses.RenderPage{
								Image:             loadPrerenderedImage("./testdata/render_testpdf_pixels_0x2000.gob", renderedPage.Image),
								PointToPixelRatio: 2.375608084404265,
							}))
							Expect(renderedPage.Image.Bounds().Size().X).To(Equal(1415))
							Expect(renderedPage.Image.Bounds().Size().Y).To(Equal(2000))
						})
					})

					Context("with both the width and height given", func() {
						Context("and the width and height being equal", func() {
							It("returns the right image, point to pixel ratio and resolution", func() {
								renderedPage, err := pdfium.RenderPageInPixels(&requests.RenderPageInPixels{
									Page:   0,
									Width:  2000,
									Height: 2000,
								})

								Expect(err).To(BeNil())
								Expect(renderedPage).To(Equal(&responses.RenderPage{
									Image:             loadPrerenderedImage("./testdata/render_testpdf_pixels_2000x2000.gob", renderedPage.Image),
									PointToPixelRatio: 2.375608084404265,
								}))

								Expect(renderedPage.Image.Bounds().Size().X).To(Equal(1415))
								Expect(renderedPage.Image.Bounds().Size().Y).To(Equal(2000))
							})
						})
						Context("and the width being larger than the height", func() {
							It("returns the right image, point to pixel ratio and resolution", func() {
								renderedPage, err := pdfium.RenderPageInPixels(&requests.RenderPageInPixels{
									Page:   0,
									Width:  4000,
									Height: 2000,
								})

								Expect(err).To(BeNil())
								Expect(renderedPage).To(Equal(&responses.RenderPage{
									Image:             loadPrerenderedImage("./testdata/render_testpdf_pixels_4000x2000.gob", renderedPage.Image),
									PointToPixelRatio: 2.375608084404265,
								}))

								Expect(renderedPage.Image.Bounds().Size().X).To(Equal(1415))
								Expect(renderedPage.Image.Bounds().Size().Y).To(Equal(2000))
							})
						})

						Context("and the height being larger than the width", func() {
							It("returns the right image, point to pixel ratio and resolution", func() {
								renderedPage, err := pdfium.RenderPageInPixels(&requests.RenderPageInPixels{
									Page:   0,
									Width:  2000,
									Height: 4000,
								})

								Expect(err).To(BeNil())
								Expect(renderedPage).To(Equal(&responses.RenderPage{
									Image:             loadPrerenderedImage("./testdata/render_testpdf_pixels_2000x4000.gob", renderedPage.Image),
									PointToPixelRatio: 3.3597884547259587,
								}))

								Expect(renderedPage.Image.Bounds().Size().X).To(Equal(2000))
								Expect(renderedPage.Image.Bounds().Size().Y).To(Equal(2829))
							})
						})
					})
				})
			})

			Context("the pages are rendered", func() {
				Context("in points", func() {
					Context("with no pages given", func() {
						It("returns an error", func() {
							renderedPage, err := pdfium.RenderPagesInDPI(&requests.RenderPagesInDPI{
								Pages: []requests.RenderPageInDPI{},
							})
							Expect(err).To(MatchError("no pages given"))
							Expect(renderedPage).To(BeNil())
						})
					})

					Context("with no DPI", func() {
						It("returns an error", func() {
							renderedPage, err := pdfium.RenderPagesInDPI(&requests.RenderPagesInDPI{
								Pages: []requests.RenderPageInDPI{
									{
										Page: 0,
									},
									{
										Page: 0,
									},
								},
							})
							Expect(err).To(MatchError("no DPI given for requested page 0"))
							Expect(renderedPage).To(BeNil())
						})
					})

					Context("with DPI 100", func() {
						It("returns the right image, point to pixel ratio and resolution", func() {
							renderedPage, err := pdfium.RenderPagesInDPI(&requests.RenderPagesInDPI{
								Pages: []requests.RenderPageInDPI{
									{
										Page: 0,
										DPI:  100,
									},
									{
										Page: 0,
										DPI:  100,
									},
								},
							})
							Expect(err).To(BeNil())
							Expect(renderedPage).To(Equal(&responses.RenderPages{
								Image: loadPrerenderedImage("./testdata/render_pages_testpdf_dpi_100.gob", renderedPage.Image),
								Pages: []responses.RenderPagesPage{
									{
										PointToPixelRatio: 1.3888888888888888,
										Width:             827,
										Height:            1170,
										X:                 0,
										Y:                 0,
									},
									{
										PointToPixelRatio: 1.3888888888888888,
										Width:             827,
										Height:            1170,
										X:                 0,
										Y:                 1170,
									},
								},
							}))
							Expect(renderedPage.Image.Bounds().Size().X).To(Equal(827))
							Expect(renderedPage.Image.Bounds().Size().Y).To(Equal(2340))
						})
					})

					Context("with DPI 300", func() {
						It("returns the right image, point to pixel ratio and resolution", func() {
							renderedPage, err := pdfium.RenderPagesInDPI(&requests.RenderPagesInDPI{
								Pages: []requests.RenderPageInDPI{
									{
										Page: 0,
										DPI:  300,
									},
									{
										Page: 0,
										DPI:  300,
									},
								},
							})
							Expect(err).To(BeNil())
							Expect(renderedPage).To(Equal(&responses.RenderPages{
								Image: loadPrerenderedImage("./testdata/render_pages_testpdf_dpi_300.gob", renderedPage.Image),
								Pages: []responses.RenderPagesPage{
									{
										PointToPixelRatio: 4.166666666666667,
										Width:             2481,
										Height:            3508,
										X:                 0,
										Y:                 0,
									},
									{
										PointToPixelRatio: 4.166666666666667,
										Width:             2481,
										Height:            3508,
										X:                 0,
										Y:                 3508,
									},
								},
							}))

							Expect(renderedPage.Image.Bounds().Size().X).To(Equal(2481))
							Expect(renderedPage.Image.Bounds().Size().Y).To(Equal(7016))
						})
					})

					Context("with different DPI per page", func() {
						It("returns the right image, point to pixel ratio and resolution", func() {
							renderedPage, err := pdfium.RenderPagesInDPI(&requests.RenderPagesInDPI{
								Pages: []requests.RenderPageInDPI{
									{
										Page: 0,
										DPI:  200,
									},
									{
										Page: 0,
										DPI:  300,
									},
								},
							})
							Expect(err).To(BeNil())
							Expect(renderedPage).To(Equal(&responses.RenderPages{
								Image: loadPrerenderedImage("./testdata/render_pages_testpdf_dpi_200_300.gob", renderedPage.Image),
								Pages: []responses.RenderPagesPage{
									{
										PointToPixelRatio: 2.7777777777777777,
										Width:             1654,
										Height:            2339,
										X:                 0,
										Y:                 0,
									},
									{
										PointToPixelRatio: 4.166666666666667,
										Width:             2481,
										Height:            3508,
										X:                 0,
										Y:                 2339,
									},
								},
							}))

							Expect(renderedPage.Image.Bounds().Size().X).To(Equal(2481))
							Expect(renderedPage.Image.Bounds().Size().Y).To(Equal(5847))
						})
					})

					Context("with padding between pages", func() {
						It("returns the right image, point to pixel ratio and resolution", func() {
							renderedPage, err := pdfium.RenderPagesInDPI(&requests.RenderPagesInDPI{
								Pages: []requests.RenderPageInDPI{
									{
										Page: 0,
										DPI:  300,
									},
									{
										Page: 0,
										DPI:  300,
									},
								},
								Padding: 50,
							})
							Expect(err).To(BeNil())
							Expect(renderedPage).To(Equal(&responses.RenderPages{
								Image: loadPrerenderedImage("./testdata/render_pages_testpdf_dpi_300_padding_50.gob", renderedPage.Image),
								Pages: []responses.RenderPagesPage{
									{
										PointToPixelRatio: 4.166666666666667,
										Width:             2481,
										Height:            3508,
										X:                 0,
										Y:                 0,
									},
									{
										PointToPixelRatio: 4.166666666666667,
										Width:             2481,
										Height:            3508,
										X:                 0,
										Y:                 3558,
									},
								},
							}))

							Expect(renderedPage.Image.Bounds().Size().X).To(Equal(2481))
							Expect(renderedPage.Image.Bounds().Size().Y).To(Equal(7066))
						})
					})
				})

				Context("in pixels", func() {
					Context("with no pages given", func() {
						It("returns an error", func() {
							renderedPage, err := pdfium.RenderPagesInPixels(&requests.RenderPagesInPixels{
								Pages: []requests.RenderPageInPixels{},
							})
							Expect(err).To(MatchError("no pages given"))
							Expect(renderedPage).To(BeNil())
						})
					})
					Context("with no width or height given", func() {
						It("returns an error", func() {
							renderedPage, err := pdfium.RenderPagesInPixels(&requests.RenderPagesInPixels{
								Pages: []requests.RenderPageInPixels{
									{
										Page: 0,
									},
									{
										Page: 0,
									},
								},
							})
							Expect(err).To(MatchError("no width or height given for requested page 0"))
							Expect(renderedPage).To(BeNil())
						})
					})

					Context("with only the width given", func() {
						It("returns the right image, point to pixel ratio and resolution", func() {
							renderedPage, err := pdfium.RenderPagesInPixels(&requests.RenderPagesInPixels{
								Pages: []requests.RenderPageInPixels{
									{
										Page:  0,
										Width: 2000,
									},
									{
										Page:  0,
										Width: 2000,
									},
								},
							})

							Expect(err).To(BeNil())
							Expect(renderedPage).To(Equal(&responses.RenderPages{
								Image: loadPrerenderedImage("./testdata/render_pages_testpdf_pixels_2000x0.gob", renderedPage.Image),
								Pages: []responses.RenderPagesPage{
									{
										PointToPixelRatio: 3.3597884547259587,
										Width:             2000,
										Height:            2829,
										X:                 0,
										Y:                 0,
									},
									{
										PointToPixelRatio: 3.3597884547259587,
										Width:             2000,
										Height:            2829,
										X:                 0,
										Y:                 2829,
									},
								},
							}))
							Expect(renderedPage.Image.Bounds().Size().X).To(Equal(2000))
							Expect(renderedPage.Image.Bounds().Size().Y).To(Equal(5658))
						})
					})

					Context("with only the height given", func() {
						It("returns the right image, point to pixel ratio and resolution", func() {
							renderedPage, err := pdfium.RenderPagesInPixels(&requests.RenderPagesInPixels{
								Pages: []requests.RenderPageInPixels{
									{
										Page:   0,
										Height: 2000,
									},
									{
										Page:   0,
										Height: 2000,
									},
								},
							})

							Expect(err).To(BeNil())
							Expect(renderedPage).To(Equal(&responses.RenderPages{
								Image: loadPrerenderedImage("./testdata/render_pages_testpdf_pixels_0x2000.gob", renderedPage.Image),
								Pages: []responses.RenderPagesPage{
									{
										PointToPixelRatio: 2.375608084404265,
										Width:             1415,
										Height:            2000,
										X:                 0,
										Y:                 0,
									},
									{
										PointToPixelRatio: 2.375608084404265,
										Width:             1415,
										Height:            2000,
										X:                 0,
										Y:                 2000,
									},
								},
							}))
							Expect(renderedPage.Image.Bounds().Size().X).To(Equal(1415))
							Expect(renderedPage.Image.Bounds().Size().Y).To(Equal(4000))
						})
					})

					Context("with both the width and height given", func() {
						Context("and the width and height being equal", func() {
							It("returns the right image, point to pixel ratio and resolution", func() {
								renderedPage, err := pdfium.RenderPagesInPixels(&requests.RenderPagesInPixels{
									Pages: []requests.RenderPageInPixels{
										{
											Page:   0,
											Width:  2000,
											Height: 2000,
										},
										{
											Page:   0,
											Width:  2000,
											Height: 2000,
										},
									},
								})

								Expect(err).To(BeNil())
								Expect(renderedPage).To(Equal(&responses.RenderPages{
									Image: loadPrerenderedImage("./testdata/render_pages_testpdf_pixels_2000x2000.gob", renderedPage.Image),
									Pages: []responses.RenderPagesPage{
										{
											PointToPixelRatio: 2.375608084404265,
											Width:             1415,
											Height:            2000,
											X:                 0,
											Y:                 0,
										},
										{
											PointToPixelRatio: 2.375608084404265,
											Width:             1415,
											Height:            2000,
											X:                 0,
											Y:                 2000,
										},
									},
								}))

								Expect(renderedPage.Image.Bounds().Size().X).To(Equal(1415))
								Expect(renderedPage.Image.Bounds().Size().Y).To(Equal(4000))
							})
						})
						Context("and the width being larger than the height", func() {
							It("returns the right image, point to pixel ratio and resolution", func() {
								renderedPage, err := pdfium.RenderPagesInPixels(&requests.RenderPagesInPixels{
									Pages: []requests.RenderPageInPixels{
										{
											Page:   0,
											Width:  4000,
											Height: 2000,
										},
										{
											Page:   0,
											Width:  4000,
											Height: 2000,
										},
									},
								})

								Expect(err).To(BeNil())
								Expect(renderedPage).To(Equal(&responses.RenderPages{
									Image: loadPrerenderedImage("./testdata/render_pages_testpdf_pixels_4000x2000.gob", renderedPage.Image),
									Pages: []responses.RenderPagesPage{
										{
											PointToPixelRatio: 2.375608084404265,
											Width:             1415,
											Height:            2000,
											X:                 0,
											Y:                 0,
										},
										{
											PointToPixelRatio: 2.375608084404265,
											Width:             1415,
											Height:            2000,
											X:                 0,
											Y:                 2000,
										},
									},
								}))

								Expect(renderedPage.Image.Bounds().Size().X).To(Equal(1415))
								Expect(renderedPage.Image.Bounds().Size().Y).To(Equal(4000))
							})
						})

						Context("and the height being larger than the width", func() {
							It("returns the right image, point to pixel ratio and resolution", func() {
								renderedPage, err := pdfium.RenderPagesInPixels(&requests.RenderPagesInPixels{
									Pages: []requests.RenderPageInPixels{
										{
											Page:   0,
											Width:  2000,
											Height: 4000,
										},
										{
											Page:   0,
											Width:  2000,
											Height: 4000,
										},
									},
								})

								Expect(err).To(BeNil())
								Expect(renderedPage).To(Equal(&responses.RenderPages{
									Image: loadPrerenderedImage("./testdata/render_pages_testpdf_pixels_2000x4000.gob", renderedPage.Image),
									Pages: []responses.RenderPagesPage{
										{
											PointToPixelRatio: 3.3597884547259587,
											Width:             2000,
											Height:            2829,
											X:                 0,
											Y:                 0,
										},
										{
											PointToPixelRatio: 3.3597884547259587,
											Width:             2000,
											Height:            2829,
											X:                 0,
											Y:                 2829,
										},
									},
								}))

								Expect(renderedPage.Image.Bounds().Size().X).To(Equal(2000))
								Expect(renderedPage.Image.Bounds().Size().Y).To(Equal(5658))
							})
						})
					})

					Context("with the width being different between pages", func() {
						It("returns the right image, point to pixel ratio and resolution", func() {
							renderedPage, err := pdfium.RenderPagesInPixels(&requests.RenderPagesInPixels{
								Pages: []requests.RenderPageInPixels{
									{
										Page:  0,
										Width: 2000,
									},
									{
										Page:  0,
										Width: 1500,
									},
								},
							})

							Expect(err).To(BeNil())
							Expect(renderedPage).To(Equal(&responses.RenderPages{
								Image: loadPrerenderedImage("./testdata/render_pages_testpdf_pixels_2000x0_1500x0.gob", renderedPage.Image),
								Pages: []responses.RenderPagesPage{
									{
										PointToPixelRatio: 3.3597884547259587,
										Width:             2000,
										Height:            2829,
										X:                 0,
										Y:                 0,
									},
									{
										PointToPixelRatio: 2.519841341044469,
										Width:             1500,
										Height:            2122,
										X:                 0,
										Y:                 2829,
									},
								},
							}))
							Expect(renderedPage.Image.Bounds().Size().X).To(Equal(2000))
							Expect(renderedPage.Image.Bounds().Size().Y).To(Equal(4951))
						})
					})

					Context("with the height being different between pages", func() {
						It("returns the right image, point to pixel ratio and resolution", func() {
							renderedPage, err := pdfium.RenderPagesInPixels(&requests.RenderPagesInPixels{
								Pages: []requests.RenderPageInPixels{
									{
										Page:   0,
										Height: 2000,
									},
									{
										Page:   0,
										Height: 1500,
									},
								},
							})

							Expect(err).To(BeNil())
							Expect(renderedPage).To(Equal(&responses.RenderPages{
								Image: loadPrerenderedImage("./testdata/render_pages_testpdf_pixels_0x2000_0x1500.gob", renderedPage.Image),
								Pages: []responses.RenderPagesPage{
									{
										PointToPixelRatio: 2.375608084404265,
										Width:             1415,
										Height:            2000,
										X:                 0,
										Y:                 0,
									},
									{
										PointToPixelRatio: 1.7817060633031987,
										Width:             1061,
										Height:            1500,
										X:                 0,
										Y:                 2000,
									},
								},
							}))
							Expect(renderedPage.Image.Bounds().Size().X).To(Equal(1415))
							Expect(renderedPage.Image.Bounds().Size().Y).To(Equal(3500))
						})
					})

					Context("with the width and height being different between pages", func() {
						It("returns the right image, point to pixel ratio and resolution", func() {
							renderedPage, err := pdfium.RenderPagesInPixels(&requests.RenderPagesInPixels{
								Pages: []requests.RenderPageInPixels{
									{
										Page:   0,
										Width:  2000,
										Height: 2000,
									},
									{
										Page:   0,
										Width:  1500,
										Height: 1500,
									},
								},
							})

							Expect(err).To(BeNil())
							Expect(renderedPage).To(Equal(&responses.RenderPages{
								Image: loadPrerenderedImage("./testdata/render_pages_testpdf_pixels_2000x2000_1500x1500.gob", renderedPage.Image),
								Pages: []responses.RenderPagesPage{
									{
										PointToPixelRatio: 2.375608084404265,
										Width:             1415,
										Height:            2000,
										X:                 0,
										Y:                 0,
									},
									{
										PointToPixelRatio: 1.7817060633031987,
										Width:             1061,
										Height:            1500,
										X:                 0,
										Y:                 2000,
									},
								},
							}))
							Expect(renderedPage.Image.Bounds().Size().X).To(Equal(1415))
							Expect(renderedPage.Image.Bounds().Size().Y).To(Equal(3500))
						})
					})

					Context("with padding between pages", func() {
						It("returns the right image, point to pixel ratio and resolution", func() {
							renderedPage, err := pdfium.RenderPagesInPixels(&requests.RenderPagesInPixels{
								Pages: []requests.RenderPageInPixels{
									{
										Page:   0,
										Width:  2000,
										Height: 2000,
									},
									{
										Page:   0,
										Width:  2000,
										Height: 2000,
									},
								},
								Padding: 50,
							})

							Expect(err).To(BeNil())
							Expect(renderedPage).To(Equal(&responses.RenderPages{
								Image: loadPrerenderedImage("./testdata/render_pages_testpdf_pixels_2000x2000_2000x2000_padding_50.gob", renderedPage.Image),
								Pages: []responses.RenderPagesPage{
									{
										PointToPixelRatio: 2.375608084404265,
										Width:             1415,
										Height:            2000,
										X:                 0,
										Y:                 0,
									},
									{
										PointToPixelRatio: 2.375608084404265,
										Width:             1415,
										Height:            2000,
										X:                 0,
										Y:                 2050,
									},
								},
							}))
							Expect(renderedPage.Image.Bounds().Size().X).To(Equal(1415))
							Expect(renderedPage.Image.Bounds().Size().Y).To(Equal(4050))
						})
					})
				})
			})
		})
	})

	// This test is only here to test the closing of an opened page.
	Context("a multipage PDF file", func() {
		BeforeEach(func() {
			pdfData, _ := ioutil.ReadFile("./testdata/test_multipage.pdf")
			pdfium.OpenDocument(&requests.OpenDocument{
				File: &pdfData,
			})
		})

		AfterEach(func() {
			pdfium.Close()
		})

		When("is opened", func() {
			Context("when another page is loaded after the first one", func() {
				Context("GetPageSize()", func() {
					It("returns the correct size for both pages", func() {
						pageSize, err := pdfium.GetPageSize(&requests.GetPageSize{
							Page: 0,
						})
						Expect(err).To(BeNil())
						Expect(pageSize).To(Equal(&responses.GetPageSize{
							Width:  595.2755737304688,
							Height: 841.8897094726562,
						}))

						pageSize, err = pdfium.GetPageSize(&requests.GetPageSize{
							Page: 1,
						})
						Expect(err).To(BeNil())
						Expect(pageSize).To(Equal(&responses.GetPageSize{
							Width:  595.2755737304688,
							Height: 841.8897094726562,
						}))
					})
				})
			})
		})
	})
})

func loadPrerenderedImage(path string, renderedImage *image.RGBA) *image.RGBA {
	err := writePrerenderedImage(path, renderedImage)
	if err != nil {
		log.Println(err)
		return nil
	}

	preRender, err := ioutil.ReadFile(path)
	if err != nil {
		return nil
	}

	buf := bytes.NewBuffer(preRender)
	dec := gob.NewDecoder(buf)

	var preRenderImage image.RGBA
	err = dec.Decode(&preRenderImage)
	return &preRenderImage
}

func writePrerenderedImage(path string, renderedImage *image.RGBA) error {
	return nil // Comment this in case of updating pdfium versions and rendering has changed.

	// Be sure to validate the difference in image to ensure rendering has not been broken.
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(renderedImage); err != nil {
		return err
	}

	if err := ioutil.WriteFile(path, buf.Bytes(), 0777); err != nil {
		return err
	}

	f, err := os.Create(path + ".png")
	if err != nil {
		return err
	}
	defer f.Close()

	err = png.Encode(f, renderedImage)
	if err != nil {
		return err
	}

	return nil
}