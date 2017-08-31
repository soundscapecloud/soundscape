package youtube

type ytdata struct {
	ResponseContext struct {
		ServiceTrackingParams []struct {
			Service string `json:"service"`
			Params  []struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			} `json:"params"`
		} `json:"serviceTrackingParams"`
		WebResponseContextExtensionData struct {
			WebResponseContextPreloadData struct {
				PreloadThumbnailUrls []string `json:"preloadThumbnailUrls"`
			} `json:"webResponseContextPreloadData"`
			YtConfigData struct {
				Csn                   string `json:"csn"`
				VisitorData           string `json:"visitorData"`
				RootVisualElementType int    `json:"rootVisualElementType"`
			} `json:"ytConfigData"`
			FeedbackDialog struct {
				PolymerOptOutFeedbackDialogRenderer struct {
					Title struct {
						Runs []struct {
							Text string `json:"text"`
						} `json:"runs"`
					} `json:"title"`
					Subtitle struct {
						Runs []struct {
							Text               string `json:"text"`
							NavigationEndpoint struct {
								URLEndpoint struct {
									URL string `json:"url"`
								} `json:"urlEndpoint"`
								WebNavigationEndpointData struct {
									URL string `json:"url"`
								} `json:"webNavigationEndpointData"`
							} `json:"navigationEndpoint,omitempty"`
						} `json:"runs"`
					} `json:"subtitle"`
					Options []struct {
						PolymerOptOutFeedbackOptionRenderer struct {
							OptionKey   string `json:"optionKey"`
							Description struct {
								Runs []struct {
									Text string `json:"text"`
								} `json:"runs"`
							} `json:"description"`
							ResponsePlaceholder struct {
								Runs []struct {
									Text string `json:"text"`
								} `json:"runs"`
							} `json:"responsePlaceholder"`
						} `json:"polymerOptOutFeedbackOptionRenderer,omitempty"`
						PolymerOptOutFeedbackNullOptionRenderer struct {
							Description struct {
								Runs []struct {
									Text string `json:"text"`
								} `json:"runs"`
							} `json:"description"`
						} `json:"polymerOptOutFeedbackNullOptionRenderer,omitempty"`
					} `json:"options"`
					Disclaimer struct {
						Runs []struct {
							Text               string `json:"text"`
							NavigationEndpoint struct {
								URLEndpoint struct {
									URL string `json:"url"`
								} `json:"urlEndpoint"`
								WebNavigationEndpointData struct {
									URL string `json:"url"`
								} `json:"webNavigationEndpointData"`
							} `json:"navigationEndpoint,omitempty"`
						} `json:"runs"`
					} `json:"disclaimer"`
					DismissButton struct {
						ButtonRenderer struct {
							Style      string `json:"style"`
							Size       string `json:"size"`
							IsDisabled bool   `json:"isDisabled"`
							Text       struct {
								Runs []struct {
									Text string `json:"text"`
								} `json:"runs"`
							} `json:"text"`
						} `json:"buttonRenderer"`
					} `json:"dismissButton"`
					SubmitButton struct {
						ButtonRenderer struct {
							Style      string `json:"style"`
							Size       string `json:"size"`
							IsDisabled bool   `json:"isDisabled"`
							Text       struct {
								Runs []struct {
									Text string `json:"text"`
								} `json:"runs"`
							} `json:"text"`
						} `json:"buttonRenderer"`
					} `json:"submitButton"`
					CloseButton struct {
						ButtonRenderer struct {
							Style      string `json:"style"`
							Size       string `json:"size"`
							IsDisabled bool   `json:"isDisabled"`
							Icon       struct {
								IconType string `json:"iconType"`
							} `json:"icon"`
						} `json:"buttonRenderer"`
					} `json:"closeButton"`
					CancelButton struct {
						ButtonRenderer struct {
							Style      string `json:"style"`
							Size       string `json:"size"`
							IsDisabled bool   `json:"isDisabled"`
							Text       struct {
								Runs []struct {
									Text string `json:"text"`
								} `json:"runs"`
							} `json:"text"`
						} `json:"buttonRenderer"`
					} `json:"cancelButton"`
				} `json:"polymerOptOutFeedbackDialogRenderer"`
			} `json:"feedbackDialog"`
		} `json:"webResponseContextExtensionData"`
	} `json:"responseContext"`
	EstimatedResults string `json:"estimatedResults"`
	Contents         struct {
		TwoColumnSearchResultsRenderer struct {
			PrimaryContents struct {
				SectionListRenderer struct {
					Contents []struct {
						ItemSectionRenderer struct {
							Contents []struct {
								VideoRenderer struct {
									VideoID   string `json:"videoId"`
									Thumbnail struct {
										Thumbnails []struct {
											URL    string `json:"url"`
											Width  int    `json:"width"`
											Height int    `json:"height"`
										} `json:"thumbnails"`
										WebThumbnailDetailsExtensionData struct {
											IsPreloaded bool `json:"isPreloaded"`
										} `json:"webThumbnailDetailsExtensionData"`
									} `json:"thumbnail"`
									Title struct {
										Accessibility struct {
											AccessibilityData struct {
												Label string `json:"label"`
											} `json:"accessibilityData"`
										} `json:"accessibility"`
										SimpleText string `json:"simpleText"`
									} `json:"title"`
									DescriptionSnippet struct {
										Runs []struct {
											Text string `json:"text"`
											Bold bool   `json:"bold,omitempty"`
										} `json:"runs"`
									} `json:"descriptionSnippet"`
									LongBylineText struct {
										Runs []struct {
											Text               string `json:"text"`
											NavigationEndpoint struct {
												ClickTrackingParams string `json:"clickTrackingParams"`
												BrowseEndpoint      struct {
													BrowseID         string `json:"browseId"`
													CanonicalBaseURL string `json:"canonicalBaseUrl"`
												} `json:"browseEndpoint"`
												WebNavigationEndpointData struct {
													URL         string `json:"url"`
													WebPageType string `json:"webPageType"`
												} `json:"webNavigationEndpointData"`
											} `json:"navigationEndpoint"`
										} `json:"runs"`
									} `json:"longBylineText"`
									PublishedTimeText struct {
										SimpleText string `json:"simpleText"`
									} `json:"publishedTimeText"`
									LengthText struct {
										Accessibility struct {
											AccessibilityData struct {
												Label string `json:"label"`
											} `json:"accessibilityData"`
										} `json:"accessibility"`
										SimpleText string `json:"simpleText"`
									} `json:"lengthText"`
									ViewCountText struct {
										SimpleText string `json:"simpleText"`
									} `json:"viewCountText"`
									NavigationEndpoint struct {
										ClickTrackingParams string `json:"clickTrackingParams"`
										WatchEndpoint       struct {
											VideoID string `json:"videoId"`
										} `json:"watchEndpoint"`
										WebNavigationEndpointData struct {
											URL         string `json:"url"`
											WebPageType string `json:"webPageType"`
										} `json:"webNavigationEndpointData"`
									} `json:"navigationEndpoint"`
									Badges []struct {
										MetadataBadgeRenderer struct {
											Style          string `json:"style"`
											Label          string `json:"label"`
											TrackingParams string `json:"trackingParams"`
										} `json:"metadataBadgeRenderer"`
									} `json:"badges"`
									OwnerBadges []struct {
										MetadataBadgeRenderer struct {
											Icon struct {
												IconType string `json:"iconType"`
											} `json:"icon"`
											Style          string `json:"style"`
											Tooltip        string `json:"tooltip"`
											TrackingParams string `json:"trackingParams"`
										} `json:"metadataBadgeRenderer"`
									} `json:"ownerBadges"`
									OwnerText struct {
										Runs []struct {
											Text               string `json:"text"`
											NavigationEndpoint struct {
												ClickTrackingParams string `json:"clickTrackingParams"`
												BrowseEndpoint      struct {
													BrowseID         string `json:"browseId"`
													CanonicalBaseURL string `json:"canonicalBaseUrl"`
												} `json:"browseEndpoint"`
												WebNavigationEndpointData struct {
													URL         string `json:"url"`
													WebPageType string `json:"webPageType"`
												} `json:"webNavigationEndpointData"`
											} `json:"navigationEndpoint"`
										} `json:"runs"`
									} `json:"ownerText"`
									ShortBylineText struct {
										Runs []struct {
											Text               string `json:"text"`
											NavigationEndpoint struct {
												ClickTrackingParams string `json:"clickTrackingParams"`
												BrowseEndpoint      struct {
													BrowseID         string `json:"browseId"`
													CanonicalBaseURL string `json:"canonicalBaseUrl"`
												} `json:"browseEndpoint"`
												WebNavigationEndpointData struct {
													URL         string `json:"url"`
													WebPageType string `json:"webPageType"`
												} `json:"webNavigationEndpointData"`
											} `json:"navigationEndpoint"`
										} `json:"runs"`
									} `json:"shortBylineText"`
									ChannelThumbnail struct {
										Thumbnails []struct {
											URL    string `json:"url"`
											Width  int    `json:"width"`
											Height int    `json:"height"`
										} `json:"thumbnails"`
									} `json:"channelThumbnail"`
									TrackingParams     string `json:"trackingParams"`
									ShowActionMenu     bool   `json:"showActionMenu"`
									ShortViewCountText struct {
										SimpleText string `json:"simpleText"`
									} `json:"shortViewCountText"`
									ThumbnailOverlays []struct {
										ThumbnailOverlayTimeStatusRenderer struct {
											Text struct {
												Accessibility struct {
													AccessibilityData struct {
														Label string `json:"label"`
													} `json:"accessibilityData"`
												} `json:"accessibility"`
												SimpleText string `json:"simpleText"`
											} `json:"text"`
											Style string `json:"style"`
										} `json:"thumbnailOverlayTimeStatusRenderer,omitempty"`
										ThumbnailOverlayToggleButtonRenderer struct {
											IsToggled     bool `json:"isToggled"`
											UntoggledIcon struct {
												IconType string `json:"iconType"`
											} `json:"untoggledIcon"`
											ToggledIcon struct {
												IconType string `json:"iconType"`
											} `json:"toggledIcon"`
											UntoggledTooltip         string `json:"untoggledTooltip"`
											ToggledTooltip           string `json:"toggledTooltip"`
											UntoggledServiceEndpoint struct {
												ClickTrackingParams  string `json:"clickTrackingParams"`
												PlaylistEditEndpoint struct {
													PlaylistID string `json:"playlistId"`
													Actions    []struct {
														AddedVideoID string `json:"addedVideoId"`
														Action       string `json:"action"`
													} `json:"actions"`
												} `json:"playlistEditEndpoint"`
											} `json:"untoggledServiceEndpoint"`
											ToggledServiceEndpoint struct {
												ClickTrackingParams  string `json:"clickTrackingParams"`
												PlaylistEditEndpoint struct {
													PlaylistID string `json:"playlistId"`
													Actions    []struct {
														Action         string `json:"action"`
														RemovedVideoID string `json:"removedVideoId"`
													} `json:"actions"`
												} `json:"playlistEditEndpoint"`
											} `json:"toggledServiceEndpoint"`
											UntoggledAccessibility struct {
												AccessibilityData struct {
													Label string `json:"label"`
												} `json:"accessibilityData"`
											} `json:"untoggledAccessibility"`
											ToggledAccessibility struct {
												AccessibilityData struct {
													Label string `json:"label"`
												} `json:"accessibilityData"`
											} `json:"toggledAccessibility"`
										} `json:"thumbnailOverlayToggleButtonRenderer,omitempty"`
									} `json:"thumbnailOverlays"`
									RichThumbnail struct {
										MovingThumbnailRenderer struct {
											MovingThumbnailDetails struct {
												Thumbnails []struct {
													URL    string `json:"url"`
													Width  int    `json:"width"`
													Height int    `json:"height"`
												} `json:"thumbnails"`
												LogAsMovingThumbnail bool `json:"logAsMovingThumbnail"`
											} `json:"movingThumbnailDetails"`
											EnableOverlay bool `json:"enableOverlay"`
										} `json:"movingThumbnailRenderer"`
									} `json:"richThumbnail"`
								} `json:"videoRenderer"`
							} `json:"contents"`
							Continuations []struct {
								NextContinuationData struct {
									Continuation        string `json:"continuation"`
									ClickTrackingParams string `json:"clickTrackingParams"`
									Label               struct {
										Runs []struct {
											Text string `json:"text"`
										} `json:"runs"`
									} `json:"label"`
								} `json:"nextContinuationData"`
							} `json:"continuations"`
							TrackingParams string `json:"trackingParams"`
						} `json:"itemSectionRenderer"`
					} `json:"contents"`
					TrackingParams string `json:"trackingParams"`
					SubMenu        struct {
						SearchSubMenuRenderer struct {
							Title struct {
								Runs []struct {
									Text string `json:"text"`
								} `json:"runs"`
							} `json:"title"`
							Groups []struct {
								SearchFilterGroupRenderer struct {
									Title struct {
										SimpleText string `json:"simpleText"`
									} `json:"title"`
									Filters []struct {
										SearchFilterRenderer struct {
											Label struct {
												SimpleText string `json:"simpleText"`
											} `json:"label"`
											NavigationEndpoint struct {
												ClickTrackingParams string `json:"clickTrackingParams"`
												SearchEndpoint      struct {
													Query  string `json:"query"`
													Params string `json:"params"`
												} `json:"searchEndpoint"`
												WebNavigationEndpointData struct {
													URL         string `json:"url"`
													WebPageType string `json:"webPageType"`
												} `json:"webNavigationEndpointData"`
											} `json:"navigationEndpoint"`
											Tooltip        string `json:"tooltip"`
											TrackingParams string `json:"trackingParams"`
										} `json:"searchFilterRenderer"`
									} `json:"filters"`
									TrackingParams string `json:"trackingParams"`
								} `json:"searchFilterGroupRenderer"`
							} `json:"groups"`
							TrackingParams string `json:"trackingParams"`
							ResultCount    struct {
								Runs []struct {
									Text string `json:"text"`
								} `json:"runs"`
							} `json:"resultCount"`
							Button struct {
								ToggleButtonRenderer struct {
									Style struct {
										StyleType string `json:"styleType"`
									} `json:"style"`
									IsToggled   bool `json:"isToggled"`
									IsDisabled  bool `json:"isDisabled"`
									DefaultIcon struct {
										IconType string `json:"iconType"`
									} `json:"defaultIcon"`
									DefaultText struct {
										SimpleText string `json:"simpleText"`
									} `json:"defaultText"`
									Accessibility struct {
										Label string `json:"label"`
									} `json:"accessibility"`
									TrackingParams string `json:"trackingParams"`
									DefaultTooltip string `json:"defaultTooltip"`
									ToggledTooltip string `json:"toggledTooltip"`
									ToggledStyle   struct {
										StyleType string `json:"styleType"`
									} `json:"toggledStyle"`
								} `json:"toggleButtonRenderer"`
							} `json:"button"`
						} `json:"searchSubMenuRenderer"`
					} `json:"subMenu"`
				} `json:"sectionListRenderer"`
			} `json:"primaryContents"`
			SecondaryContents struct {
				SecondarySearchContainerRenderer struct {
					Contents []struct {
						SearchMpuAdRenderer struct {
							VerticalIds    []string `json:"verticalIds"`
							TrackingParams string   `json:"trackingParams"`
						} `json:"searchMpuAdRenderer"`
					} `json:"contents"`
					TrackingParams string `json:"trackingParams"`
				} `json:"secondarySearchContainerRenderer"`
			} `json:"secondaryContents"`
		} `json:"twoColumnSearchResultsRenderer"`
	} `json:"contents"`
	TrackingParams string `json:"trackingParams"`
	Topbar         struct {
		DesktopTopbarRenderer struct {
			Logo struct {
				TopbarLogoRenderer struct {
					IconImage struct {
						IconType string `json:"iconType"`
					} `json:"iconImage"`
					TooltipText struct {
						Runs []struct {
							Text string `json:"text"`
						} `json:"runs"`
					} `json:"tooltipText"`
					Endpoint struct {
						ClickTrackingParams string `json:"clickTrackingParams"`
						BrowseEndpoint      struct {
							BrowseID string `json:"browseId"`
						} `json:"browseEndpoint"`
						WebNavigationEndpointData struct {
							URL         string `json:"url"`
							WebPageType string `json:"webPageType"`
						} `json:"webNavigationEndpointData"`
					} `json:"endpoint"`
					TrackingParams string `json:"trackingParams"`
				} `json:"topbarLogoRenderer"`
			} `json:"logo"`
			Searchbox struct {
				FusionSearchboxRenderer struct {
					Icon struct {
						IconType string `json:"iconType"`
					} `json:"icon"`
					PlaceholderText struct {
						Runs []struct {
							Text string `json:"text"`
						} `json:"runs"`
					} `json:"placeholderText"`
					Config struct {
						WebSearchboxConfig struct {
							RequestLanguage     string `json:"requestLanguage"`
							RequestDomain       string `json:"requestDomain"`
							HasOnscreenKeyboard bool   `json:"hasOnscreenKeyboard"`
							FocusSearchbox      bool   `json:"focusSearchbox"`
						} `json:"webSearchboxConfig"`
					} `json:"config"`
					TrackingParams string `json:"trackingParams"`
					SearchEndpoint struct {
						ClickTrackingParams string `json:"clickTrackingParams"`
						SearchEndpoint      struct {
							Query string `json:"query"`
						} `json:"searchEndpoint"`
						WebNavigationEndpointData struct {
							URL         string `json:"url"`
							WebPageType string `json:"webPageType"`
						} `json:"webNavigationEndpointData"`
					} `json:"searchEndpoint"`
				} `json:"fusionSearchboxRenderer"`
			} `json:"searchbox"`
			TrackingParams string `json:"trackingParams"`
			TopbarButtons  []struct {
				ButtonRenderer struct {
					Style      string `json:"style"`
					Size       string `json:"size"`
					IsDisabled bool   `json:"isDisabled"`
					Icon       struct {
						IconType string `json:"iconType"`
					} `json:"icon"`
					NavigationEndpoint struct {
						ClickTrackingParams string `json:"clickTrackingParams"`
						SignInEndpoint      struct {
							NextEndpoint struct {
								ClickTrackingParams string `json:"clickTrackingParams"`
								UploadEndpoint      struct {
									Hack bool `json:"hack"`
								} `json:"uploadEndpoint"`
								WebNavigationEndpointData struct {
									URL string `json:"url"`
								} `json:"webNavigationEndpointData"`
							} `json:"nextEndpoint"`
						} `json:"signInEndpoint"`
						WebNavigationEndpointData struct {
							URL string `json:"url"`
						} `json:"webNavigationEndpointData"`
					} `json:"navigationEndpoint"`
					Accessibility struct {
						Label string `json:"label"`
					} `json:"accessibility"`
					Tooltip           string `json:"tooltip"`
					TrackingParams    string `json:"trackingParams"`
					AccessibilityData struct {
						AccessibilityData struct {
							Label string `json:"label"`
						} `json:"accessibilityData"`
					} `json:"accessibilityData"`
				} `json:"buttonRenderer,omitempty"`
				TopbarMenuButtonRenderer struct {
					Icon struct {
						IconType string `json:"iconType"`
					} `json:"icon"`
					MenuRenderer struct {
						MultiPageMenuRenderer struct {
							Sections []struct {
								MultiPageMenuSectionRenderer struct {
									Items []struct {
										CompactLinkRenderer struct {
											Icon struct {
												IconType string `json:"iconType"`
											} `json:"icon"`
											Title struct {
												Runs []struct {
													Text string `json:"text"`
												} `json:"runs"`
											} `json:"title"`
											NavigationEndpoint struct {
												ClickTrackingParams string `json:"clickTrackingParams"`
												URLEndpoint         struct {
													URL    string `json:"url"`
													Target string `json:"target"`
												} `json:"urlEndpoint"`
												WebNavigationEndpointData struct {
													URL string `json:"url"`
												} `json:"webNavigationEndpointData"`
											} `json:"navigationEndpoint"`
											TrackingParams string `json:"trackingParams"`
										} `json:"compactLinkRenderer"`
									} `json:"items"`
									TrackingParams string `json:"trackingParams"`
								} `json:"multiPageMenuSectionRenderer"`
							} `json:"sections"`
							TrackingParams string `json:"trackingParams"`
						} `json:"multiPageMenuRenderer"`
					} `json:"menuRenderer"`
					TrackingParams string `json:"trackingParams"`
					Accessibility  struct {
						AccessibilityData struct {
							Label string `json:"label"`
						} `json:"accessibilityData"`
					} `json:"accessibility"`
					Tooltip string `json:"tooltip"`
					Style   string `json:"style"`
				} `json:"topbarMenuButtonRenderer,omitempty"`
			} `json:"topbarButtons"`
			HotkeyDialog struct {
				HotkeyDialogRenderer struct {
					Title struct {
						Runs []struct {
							Text string `json:"text"`
						} `json:"runs"`
					} `json:"title"`
					Sections []struct {
						HotkeyDialogSectionRenderer struct {
							Title struct {
								Runs []struct {
									Text string `json:"text"`
								} `json:"runs"`
							} `json:"title"`
							Options []struct {
								HotkeyDialogSectionOptionRenderer struct {
									Label struct {
										Runs []struct {
											Text string `json:"text"`
										} `json:"runs"`
									} `json:"label"`
									Hotkey string `json:"hotkey"`
								} `json:"hotkeyDialogSectionOptionRenderer"`
							} `json:"options"`
						} `json:"hotkeyDialogSectionRenderer"`
					} `json:"sections"`
					DismissButton struct {
						ButtonRenderer struct {
							Style      string `json:"style"`
							Size       string `json:"size"`
							IsDisabled bool   `json:"isDisabled"`
							Text       struct {
								SimpleText string `json:"simpleText"`
							} `json:"text"`
							TrackingParams string `json:"trackingParams"`
						} `json:"buttonRenderer"`
					} `json:"dismissButton"`
				} `json:"hotkeyDialogRenderer"`
			} `json:"hotkeyDialog"`
		} `json:"desktopTopbarRenderer"`
	} `json:"topbar"`
	AdSafetyReason struct {
	} `json:"adSafetyReason"`
}
