# UserinfoResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Birthdate** | Pointer to **string** | End-User&#39;s birthday, represented as an ISO 8601:2004 [ISO8601‑2004] YYYY-MM-DD format. The year MAY be 0000, indicating that it is omitted. To represent only the year, YYYY format is allowed. Note that depending on the underlying platform&#39;s date related function, providing just year can result in varying month and day, so the implementers need to take this factor into account to correctly process the dates. | [optional] 
**Email** | Pointer to **string** | End-User&#39;s preferred e-mail address. Its value MUST conform to the RFC 5322 [RFC5322] addr-spec syntax. The RP MUST NOT rely upon this value being unique, as discussed in Section 5.7. | [optional] 
**EmailVerified** | Pointer to **bool** | True if the End-User&#39;s e-mail address has been verified; otherwise false. When this Claim Value is true, this means that the OP took affirmative steps to ensure that this e-mail address was controlled by the End-User at the time the verification was performed. The means by which an e-mail address is verified is context-specific, and dependent upon the trust framework or contractual agreements within which the parties are operating. | [optional] 
**FamilyName** | Pointer to **string** | Surname(s) or last name(s) of the End-User. Note that in some cultures, people can have multiple family names or no family name; all can be present, with the names being separated by space characters. | [optional] 
**Gender** | Pointer to **string** | End-User&#39;s gender. Values defined by this specification are female and male. Other values MAY be used when neither of the defined values are applicable. | [optional] 
**GivenName** | Pointer to **string** | Given name(s) or first name(s) of the End-User. Note that in some cultures, people can have multiple given names; all can be present, with the names being separated by space characters. | [optional] 
**Locale** | Pointer to **string** | End-User&#39;s locale, represented as a BCP47 [RFC5646] language tag. This is typically an ISO 639-1 Alpha-2 [ISO639‑1] language code in lowercase and an ISO 3166-1 Alpha-2 [ISO3166‑1] country code in uppercase, separated by a dash. For example, en-US or fr-CA. As a compatibility note, some implementations have used an underscore as the separator rather than a dash, for example, en_US; Relying Parties MAY choose to accept this locale syntax as well. | [optional] 
**MiddleName** | Pointer to **string** | Middle name(s) of the End-User. Note that in some cultures, people can have multiple middle names; all can be present, with the names being separated by space characters. Also note that in some cultures, middle names are not used. | [optional] 
**Name** | Pointer to **string** | End-User&#39;s full name in displayable form including all name parts, possibly including titles and suffixes, ordered according to the End-User&#39;s locale and preferences. | [optional] 
**Nickname** | Pointer to **string** | Casual name of the End-User that may or may not be the same as the given_name. For instance, a nickname value of Mike might be returned alongside a given_name value of Michael. | [optional] 
**PhoneNumber** | Pointer to **string** | End-User&#39;s preferred telephone number. E.164 [E.164] is RECOMMENDED as the format of this Claim, for example, +1 (425) 555-1212 or +56 (2) 687 2400. If the phone number contains an extension, it is RECOMMENDED that the extension be represented using the RFC 3966 [RFC3966] extension syntax, for example, +1 (604) 555-1234;ext&#x3D;5678. | [optional] 
**PhoneNumberVerified** | Pointer to **bool** | True if the End-User&#39;s phone number has been verified; otherwise false. When this Claim Value is true, this means that the OP took affirmative steps to ensure that this phone number was controlled by the End-User at the time the verification was performed. The means by which a phone number is verified is context-specific, and dependent upon the trust framework or contractual agreements within which the parties are operating. When true, the phone_number Claim MUST be in E.164 format and any extensions MUST be represented in RFC 3966 format. | [optional] 
**Picture** | Pointer to **string** | URL of the End-User&#39;s profile picture. This URL MUST refer to an image file (for example, a PNG, JPEG, or GIF image file), rather than to a Web page containing an image. Note that this URL SHOULD specifically reference a profile photo of the End-User suitable for displaying when describing the End-User, rather than an arbitrary photo taken by the End-User. | [optional] 
**PreferredUsername** | Pointer to **string** | Non-unique shorthand name by which the End-User wishes to be referred to at the RP, such as janedoe or j.doe. This value MAY be any valid JSON string including special characters such as @, /, or whitespace. | [optional] 
**Profile** | Pointer to **string** | URL of the End-User&#39;s profile page. The contents of this Web page SHOULD be about the End-User. | [optional] 
**Sub** | Pointer to **string** | Subject - Identifier for the End-User at the IssuerURL. | [optional] 
**UpdatedAt** | Pointer to **int64** | Time the End-User&#39;s information was last updated. Its value is a JSON number representing the number of seconds from 1970-01-01T0:0:0Z as measured in UTC until the date/time. | [optional] 
**Website** | Pointer to **string** | URL of the End-User&#39;s Web page or blog. This Web page SHOULD contain information published by the End-User or an organization that the End-User is affiliated with. | [optional] 
**Zoneinfo** | Pointer to **string** | String from zoneinfo [zoneinfo] time zone database representing the End-User&#39;s time zone. For example, Europe/Paris or America/Los_Angeles. | [optional] 

## Methods

### NewUserinfoResponse

`func NewUserinfoResponse() *UserinfoResponse`

NewUserinfoResponse instantiates a new UserinfoResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUserinfoResponseWithDefaults

`func NewUserinfoResponseWithDefaults() *UserinfoResponse`

NewUserinfoResponseWithDefaults instantiates a new UserinfoResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetBirthdate

`func (o *UserinfoResponse) GetBirthdate() string`

GetBirthdate returns the Birthdate field if non-nil, zero value otherwise.

### GetBirthdateOk

`func (o *UserinfoResponse) GetBirthdateOk() (*string, bool)`

GetBirthdateOk returns a tuple with the Birthdate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBirthdate

`func (o *UserinfoResponse) SetBirthdate(v string)`

SetBirthdate sets Birthdate field to given value.

### HasBirthdate

`func (o *UserinfoResponse) HasBirthdate() bool`

HasBirthdate returns a boolean if a field has been set.

### GetEmail

`func (o *UserinfoResponse) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *UserinfoResponse) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *UserinfoResponse) SetEmail(v string)`

SetEmail sets Email field to given value.

### HasEmail

`func (o *UserinfoResponse) HasEmail() bool`

HasEmail returns a boolean if a field has been set.

### GetEmailVerified

`func (o *UserinfoResponse) GetEmailVerified() bool`

GetEmailVerified returns the EmailVerified field if non-nil, zero value otherwise.

### GetEmailVerifiedOk

`func (o *UserinfoResponse) GetEmailVerifiedOk() (*bool, bool)`

GetEmailVerifiedOk returns a tuple with the EmailVerified field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmailVerified

`func (o *UserinfoResponse) SetEmailVerified(v bool)`

SetEmailVerified sets EmailVerified field to given value.

### HasEmailVerified

`func (o *UserinfoResponse) HasEmailVerified() bool`

HasEmailVerified returns a boolean if a field has been set.

### GetFamilyName

`func (o *UserinfoResponse) GetFamilyName() string`

GetFamilyName returns the FamilyName field if non-nil, zero value otherwise.

### GetFamilyNameOk

`func (o *UserinfoResponse) GetFamilyNameOk() (*string, bool)`

GetFamilyNameOk returns a tuple with the FamilyName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFamilyName

`func (o *UserinfoResponse) SetFamilyName(v string)`

SetFamilyName sets FamilyName field to given value.

### HasFamilyName

`func (o *UserinfoResponse) HasFamilyName() bool`

HasFamilyName returns a boolean if a field has been set.

### GetGender

`func (o *UserinfoResponse) GetGender() string`

GetGender returns the Gender field if non-nil, zero value otherwise.

### GetGenderOk

`func (o *UserinfoResponse) GetGenderOk() (*string, bool)`

GetGenderOk returns a tuple with the Gender field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGender

`func (o *UserinfoResponse) SetGender(v string)`

SetGender sets Gender field to given value.

### HasGender

`func (o *UserinfoResponse) HasGender() bool`

HasGender returns a boolean if a field has been set.

### GetGivenName

`func (o *UserinfoResponse) GetGivenName() string`

GetGivenName returns the GivenName field if non-nil, zero value otherwise.

### GetGivenNameOk

`func (o *UserinfoResponse) GetGivenNameOk() (*string, bool)`

GetGivenNameOk returns a tuple with the GivenName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGivenName

`func (o *UserinfoResponse) SetGivenName(v string)`

SetGivenName sets GivenName field to given value.

### HasGivenName

`func (o *UserinfoResponse) HasGivenName() bool`

HasGivenName returns a boolean if a field has been set.

### GetLocale

`func (o *UserinfoResponse) GetLocale() string`

GetLocale returns the Locale field if non-nil, zero value otherwise.

### GetLocaleOk

`func (o *UserinfoResponse) GetLocaleOk() (*string, bool)`

GetLocaleOk returns a tuple with the Locale field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLocale

`func (o *UserinfoResponse) SetLocale(v string)`

SetLocale sets Locale field to given value.

### HasLocale

`func (o *UserinfoResponse) HasLocale() bool`

HasLocale returns a boolean if a field has been set.

### GetMiddleName

`func (o *UserinfoResponse) GetMiddleName() string`

GetMiddleName returns the MiddleName field if non-nil, zero value otherwise.

### GetMiddleNameOk

`func (o *UserinfoResponse) GetMiddleNameOk() (*string, bool)`

GetMiddleNameOk returns a tuple with the MiddleName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMiddleName

`func (o *UserinfoResponse) SetMiddleName(v string)`

SetMiddleName sets MiddleName field to given value.

### HasMiddleName

`func (o *UserinfoResponse) HasMiddleName() bool`

HasMiddleName returns a boolean if a field has been set.

### GetName

`func (o *UserinfoResponse) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *UserinfoResponse) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *UserinfoResponse) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *UserinfoResponse) HasName() bool`

HasName returns a boolean if a field has been set.

### GetNickname

`func (o *UserinfoResponse) GetNickname() string`

GetNickname returns the Nickname field if non-nil, zero value otherwise.

### GetNicknameOk

`func (o *UserinfoResponse) GetNicknameOk() (*string, bool)`

GetNicknameOk returns a tuple with the Nickname field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNickname

`func (o *UserinfoResponse) SetNickname(v string)`

SetNickname sets Nickname field to given value.

### HasNickname

`func (o *UserinfoResponse) HasNickname() bool`

HasNickname returns a boolean if a field has been set.

### GetPhoneNumber

`func (o *UserinfoResponse) GetPhoneNumber() string`

GetPhoneNumber returns the PhoneNumber field if non-nil, zero value otherwise.

### GetPhoneNumberOk

`func (o *UserinfoResponse) GetPhoneNumberOk() (*string, bool)`

GetPhoneNumberOk returns a tuple with the PhoneNumber field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPhoneNumber

`func (o *UserinfoResponse) SetPhoneNumber(v string)`

SetPhoneNumber sets PhoneNumber field to given value.

### HasPhoneNumber

`func (o *UserinfoResponse) HasPhoneNumber() bool`

HasPhoneNumber returns a boolean if a field has been set.

### GetPhoneNumberVerified

`func (o *UserinfoResponse) GetPhoneNumberVerified() bool`

GetPhoneNumberVerified returns the PhoneNumberVerified field if non-nil, zero value otherwise.

### GetPhoneNumberVerifiedOk

`func (o *UserinfoResponse) GetPhoneNumberVerifiedOk() (*bool, bool)`

GetPhoneNumberVerifiedOk returns a tuple with the PhoneNumberVerified field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPhoneNumberVerified

`func (o *UserinfoResponse) SetPhoneNumberVerified(v bool)`

SetPhoneNumberVerified sets PhoneNumberVerified field to given value.

### HasPhoneNumberVerified

`func (o *UserinfoResponse) HasPhoneNumberVerified() bool`

HasPhoneNumberVerified returns a boolean if a field has been set.

### GetPicture

`func (o *UserinfoResponse) GetPicture() string`

GetPicture returns the Picture field if non-nil, zero value otherwise.

### GetPictureOk

`func (o *UserinfoResponse) GetPictureOk() (*string, bool)`

GetPictureOk returns a tuple with the Picture field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPicture

`func (o *UserinfoResponse) SetPicture(v string)`

SetPicture sets Picture field to given value.

### HasPicture

`func (o *UserinfoResponse) HasPicture() bool`

HasPicture returns a boolean if a field has been set.

### GetPreferredUsername

`func (o *UserinfoResponse) GetPreferredUsername() string`

GetPreferredUsername returns the PreferredUsername field if non-nil, zero value otherwise.

### GetPreferredUsernameOk

`func (o *UserinfoResponse) GetPreferredUsernameOk() (*string, bool)`

GetPreferredUsernameOk returns a tuple with the PreferredUsername field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPreferredUsername

`func (o *UserinfoResponse) SetPreferredUsername(v string)`

SetPreferredUsername sets PreferredUsername field to given value.

### HasPreferredUsername

`func (o *UserinfoResponse) HasPreferredUsername() bool`

HasPreferredUsername returns a boolean if a field has been set.

### GetProfile

`func (o *UserinfoResponse) GetProfile() string`

GetProfile returns the Profile field if non-nil, zero value otherwise.

### GetProfileOk

`func (o *UserinfoResponse) GetProfileOk() (*string, bool)`

GetProfileOk returns a tuple with the Profile field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProfile

`func (o *UserinfoResponse) SetProfile(v string)`

SetProfile sets Profile field to given value.

### HasProfile

`func (o *UserinfoResponse) HasProfile() bool`

HasProfile returns a boolean if a field has been set.

### GetSub

`func (o *UserinfoResponse) GetSub() string`

GetSub returns the Sub field if non-nil, zero value otherwise.

### GetSubOk

`func (o *UserinfoResponse) GetSubOk() (*string, bool)`

GetSubOk returns a tuple with the Sub field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSub

`func (o *UserinfoResponse) SetSub(v string)`

SetSub sets Sub field to given value.

### HasSub

`func (o *UserinfoResponse) HasSub() bool`

HasSub returns a boolean if a field has been set.

### GetUpdatedAt

`func (o *UserinfoResponse) GetUpdatedAt() int64`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *UserinfoResponse) GetUpdatedAtOk() (*int64, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *UserinfoResponse) SetUpdatedAt(v int64)`

SetUpdatedAt sets UpdatedAt field to given value.

### HasUpdatedAt

`func (o *UserinfoResponse) HasUpdatedAt() bool`

HasUpdatedAt returns a boolean if a field has been set.

### GetWebsite

`func (o *UserinfoResponse) GetWebsite() string`

GetWebsite returns the Website field if non-nil, zero value otherwise.

### GetWebsiteOk

`func (o *UserinfoResponse) GetWebsiteOk() (*string, bool)`

GetWebsiteOk returns a tuple with the Website field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWebsite

`func (o *UserinfoResponse) SetWebsite(v string)`

SetWebsite sets Website field to given value.

### HasWebsite

`func (o *UserinfoResponse) HasWebsite() bool`

HasWebsite returns a boolean if a field has been set.

### GetZoneinfo

`func (o *UserinfoResponse) GetZoneinfo() string`

GetZoneinfo returns the Zoneinfo field if non-nil, zero value otherwise.

### GetZoneinfoOk

`func (o *UserinfoResponse) GetZoneinfoOk() (*string, bool)`

GetZoneinfoOk returns a tuple with the Zoneinfo field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetZoneinfo

`func (o *UserinfoResponse) SetZoneinfo(v string)`

SetZoneinfo sets Zoneinfo field to given value.

### HasZoneinfo

`func (o *UserinfoResponse) HasZoneinfo() bool`

HasZoneinfo returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


