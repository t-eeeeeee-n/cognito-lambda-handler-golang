package cognito

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"
)

// 定数
const (
	nHex = "FFFFFFFFFFFFFFFFC90FDAA22168C234C4C6628B80DC1CD1" +
		"29024E088A67CC74020BBEA63B139B22514A08798E3404DD" +
		"EF9519B3CD3A431B302B0A6DF25F14374FE1356D6D51C245" +
		"E485B576625E7EC6F44C42E9A637ED6B0BFF5CB6F406B7ED" +
		"EE386BFB5A899FA5AE9F24117C4B1FE649286651ECE45B3D" +
		"C2007CB8A163BF0598DA48361C55D39A69163FA8FD24CF5F" +
		"83655D23DCA3AD961C62F356208552BB9ED529077096966D" +
		"670C354E4ABC9804F1746C08CA18217C32905E462E36CE3B" +
		"E39E772C180E86039B2783A2EC07A28FB5C55DF06F4C52C9" +
		"DE2BCBF6955817183995497CEA956AE515D2261898FA0510" +
		"15728E5A8AAAC42DAD33170D04507A33A85521ABDF1CBA64" +
		"ECFB850458DBEF0A8AEA71575D060C7DB3970F85A6E1E4C7" +
		"ABF5AE8CDB0933D71E8C94E04A25619DCEE3D2261AD2EE6B" +
		"F12FFA06D98A0864D87602733EC86A64521F2B18177B200C" +
		"BBE117577A615D6C770988C0BAD946E208E24FA074E5AB31" +
		"43DB5BFCE0FD108E4B82D120A93AD2CAFFFFFFFFFFFFFFFF"
	gHex     = "2" // 小さい値（通常2）
	infoBits = "Caldera Derived Key"
)

// SRP はSRP認証のための構造体
type SRP struct {
	Username     string
	Password     string
	PoolId       string
	PoolName     string
	ClientId     string
	ClientSecret string
	BigN         *big.Int
	G            *big.Int
	K            *big.Int
	A            *big.Int
	BigA         *big.Int
}

// NewCognitoSRP は新しいSRPオブジェクトを作成
func NewCognitoSRP(username, password, poolId, clientId string, clientSecret string) (*SRP, error) {
	bigN, err := hexToBig(nHex)
	if err != nil {
		return nil, fmt.Errorf("failed to convert nHex to big.Int: %w", err)
	}

	g, err := hexToBig(gHex)
	if err != nil {
		return nil, fmt.Errorf("failed to convert gHex to big.Int: %w", err)
	}

	c := &SRP{
		Username:     username,
		Password:     password,
		PoolId:       poolId,
		ClientId:     clientId,
		ClientSecret: clientSecret,
		BigN:         bigN,
		G:            g,
	}

	if !strings.Contains(poolId, "_") {
		return nil, fmt.Errorf("invalid Cognito User Pool ID (%s), must be in format: '<region>_<pool name>'", poolId)
	}
	c.PoolName = strings.Split(poolId, "_")[1]

	// k値の計算
	c.K, err = hexToBig(hexHash("00" + nHex + "0" + gHex))
	if err != nil {
		return nil, fmt.Errorf("failed to calculate k value: %w", err)
	}

	// ランダムなa値を生成
	c.A, err = c.generateRandomSmallA()
	if err != nil {
		return nil, fmt.Errorf("failed to generate random a value: %w", err)
	}

	// A値の計算
	c.BigA, err = c.calculateA()
	if err != nil {
		return nil, fmt.Errorf("failed to calculate A value: %w", err)
	}

	return c, nil
}

// GetUsername は設定されたCognitoユーザー名を返却
func (csrp *SRP) GetUsername() string {
	return csrp.Username
}

// GetClientId は設定されたCognitoクライアントIDを返却
func (csrp *SRP) GetClientId() string {
	return csrp.ClientId
}

// GetUserPoolId は設定されたCognitoユーザープールIDを返却
func (csrp *SRP) GetUserPoolId() string {
	return csrp.PoolId
}

// GetUserPoolName は設定されたCognitoユーザープール名を返却
func (csrp *SRP) GetUserPoolName() string {
	return csrp.PoolName
}

// GetAuthParams はInitiateAuthリクエストに必要な認証パラメータを返却
func (csrp *SRP) GetAuthParams() map[string]string {
	params := map[string]string{
		"USERNAME": csrp.Username,
		"SRP_A":    bigToHex(csrp.BigA),
	}

	secret, err := csrp.GetSecretHash(csrp.Username)
	if err != nil {
		fmt.Printf("Failed to generate SECRET_HASH: %v\n", err)
	} else {
		params["SECRET_HASH"] = secret
	}

	return params
}

// GetSecretHash は、クライアントがシークレットで構成されている場合に必要なシークレットハッシュを生成
func (csrp *SRP) GetSecretHash(username string) (string, error) {
	var (
		msg = username + csrp.ClientId
		key = []byte(csrp.ClientSecret)
		h   = hmac.New(sha256.New, key)
	)

	h.Write([]byte(msg))

	sh := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return sh, nil
}

// PasswordVerifierChallenge はPASSWORD_VERIFIERチャレンジを完了するために使用するChallengeResponsesを返却
func (csrp *SRP) PasswordVerifierChallenge(challengeParms map[string]string, ts time.Time) (map[string]string, error) {
	var (
		internalUsername = challengeParms["USERNAME"]
		userId           = challengeParms["USER_ID_FOR_SRP"]
		saltHex          = challengeParms["SALT"]
		srpBHex          = challengeParms["SRP_B"]
		secretBlockB64   = challengeParms["SECRET_BLOCK"]

		timestamp = ts.In(time.UTC).Format("Mon Jan 2 03:04:05 MST 2006")
	)

	// srpBHexの変換とエラーチェック
	srpB, err := hexToBig(srpBHex)
	if err != nil {
		return nil, fmt.Errorf("failed to convert SRP_B to big.Int: %w", err)
	}

	// saltHexの変換とエラーチェック
	salt, err := hexToBig(saltHex)
	if err != nil {
		return nil, fmt.Errorf("failed to convert SALT to big.Int: %w", err)
	}

	// パスワード認証キーの取得
	hkdf, err := csrp.getPasswordAuthenticationKey(userId, csrp.Password, srpB, salt)
	if err != nil {
		return nil, fmt.Errorf("failed to get password authentication key: %w", err)
	}

	secretBlockBytes, err := base64.StdEncoding.DecodeString(secretBlockB64)
	if err != nil {
		return nil, fmt.Errorf("unable to decode challenge parameter 'SECRET_BLOCK', %s", err.Error())
	}

	msg := csrp.PoolName + userId + string(secretBlockBytes) + timestamp
	hmacObj := hmac.New(sha256.New, hkdf)
	hmacObj.Write([]byte(msg))
	signature := base64.StdEncoding.EncodeToString(hmacObj.Sum(nil))

	response := map[string]string{
		"TIMESTAMP":                   timestamp,
		"USERNAME":                    internalUsername,
		"PASSWORD_CLAIM_SECRET_BLOCK": secretBlockB64,
		"PASSWORD_CLAIM_SIGNATURE":    signature,
	}
	if secret, err := csrp.GetSecretHash(internalUsername); err == nil {
		response["SECRET_HASH"] = secret
	}

	return response, nil
}

// generateRandomSmallA はランダムなa値を生成
func (csrp *SRP) generateRandomSmallA() (*big.Int, error) {
	randomLongInt, err := getRandom(128)
	if err != nil {
		return nil, fmt.Errorf("failed to generate small a value: %w", err)
	}
	return big.NewInt(0).Mod(randomLongInt, csrp.BigN), nil
}

// calculateA はA値を計算
func (csrp *SRP) calculateA() (*big.Int, error) {
	bigA := big.NewInt(0).Exp(csrp.G, csrp.A, csrp.BigN)
	if big.NewInt(0).Mod(bigA, csrp.BigN).Cmp(big.NewInt(0)) == 0 {
		return nil, fmt.Errorf("safety check for A failed: A must not be divisible by N")
	}
	return bigA, nil
}

// getPasswordAuthenticationKey SRPプロトコルにおけるパスワード認証キーを生成
func (csrp *SRP) getPasswordAuthenticationKey(username, password string, bigB, salt *big.Int) ([]byte, error) {
	// ユーザー名、パスワード、プール名を組み合わせてユーザーパスハッシュを作成
	userPass := fmt.Sprintf("%s%s:%s", csrp.PoolName, username, password)
	userPassHash := hashSha256([]byte(userPass))

	// U値の計算
	uVal, err := calculateU(csrp.BigA, bigB)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate U value: %w", err)
	}

	// x値の計算
	xVal, err := hexToBig(hexHash(padHex(salt.Text(16)) + userPassHash))
	if err != nil {
		return nil, fmt.Errorf("error converting hex to big.Int: %w", err)
	}

	// g^x mod N の計算
	gModPowXN := big.NewInt(0).Exp(csrp.G, xVal, csrp.BigN)

	// S = (B - k * g^x) ^ (a + u * x) mod N の計算
	intVal1 := big.NewInt(0).Sub(bigB, big.NewInt(0).Mul(csrp.K, gModPowXN))
	intVal2 := big.NewInt(0).Add(csrp.A, big.NewInt(0).Mul(uVal, xVal))
	sVal := big.NewInt(0).Exp(intVal1, intVal2, csrp.BigN)

	return computeHKDF(padHex(sVal.Text(16)), padHex(bigToHex(uVal))), nil
}

// hashSha256 はSHA256ハッシュを返す
func hashSha256(buf []byte) string {
	a := sha256.New()
	a.Write(buf)
	return hex.EncodeToString(a.Sum(nil))
}

// hexHash はSHA256ハッシュを16進文字列として返す
func hexHash(hexStr string) string {
	buf, _ := hex.DecodeString(hexStr)
	return hashSha256(buf)
}

// hexToBig は16進数文字列を *big.Int に変換
func hexToBig(hexStr string) (*big.Int, error) {
	i, ok := big.NewInt(0).SetString(hexStr, 16)
	if !ok {
		return nil, fmt.Errorf("unable to convert \"%s\" to big Int", hexStr)
	}
	return i, nil
}

// bigToHex *big.Int を16進数文字列に変換して返す
func bigToHex(val *big.Int) string {
	return val.Text(16)
}

// getRandom は指定したバイト数のランダムな値を生成
func getRandom(n int) (*big.Int, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random bytes: %w", err)
	}
	randomBigInt, err := hexToBig(hex.EncodeToString(b))
	if err != nil {
		return nil, fmt.Errorf("failed to convert random bytes to big.Int: %w", err)
	}
	return randomBigInt, nil
}

// padHex 16進文字列を必要に応じてゼロでパディング
func padHex(hexStr string) string {
	if len(hexStr)%2 == 1 {
		hexStr = fmt.Sprintf("0%s", hexStr)
	} else if strings.Contains("89ABCDEFabcdef", string(hexStr[0])) {
		hexStr = fmt.Sprintf("00%s", hexStr)
	}

	return hexStr
}

// computeHKDF HKDF (HMAC-based Extract-and-Expand Key Derivation Function) を用いてキーを生成
func computeHKDF(ikm, salt string) []byte {
	ikmb, _ := hex.DecodeString(ikm)
	saltb, _ := hex.DecodeString(salt)

	extractor := hmac.New(sha256.New, saltb)
	extractor.Write(ikmb)
	prk := extractor.Sum(nil)
	infoBitsUpdate := append([]byte(infoBits), byte(1))
	extractor = hmac.New(sha256.New, prk)
	extractor.Write(infoBitsUpdate)
	hmacHash := extractor.Sum(nil)

	return hmacHash[:16]
}

// calculateU SRPプロトコルにおけるハッシュ値 u を計算
func calculateU(bigA, bigB *big.Int) (*big.Int, error) {

	hexResult := hexHash(padHex(bigA.Text(16)) + padHex(bigB.Text(16)))
	u, err := hexToBig(hexResult)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate u value: %w", err)
	}
	return u, nil
}
