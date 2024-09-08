export const getFlagByCountry = (countryCode: string) => {
	if (!countryCode || countryCode.length !== 2 || !/^[a-zA-Z]+$/.test(countryCode)) {
		return "🏳️";
	}

	const code = countryCode.toUpperCase();
	const offset = 127397;
	const flag = Array.from(code)
		.map((letter) => String.fromCodePoint(letter.charCodeAt(0) + offset))
		.join("");

	return flag;
};
