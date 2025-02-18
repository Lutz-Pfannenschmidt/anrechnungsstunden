export function isValidEmail(email) {
	const atSymbolIndex = email.indexOf("@");
	return (
		atSymbolIndex > 0 &&
		email.indexOf(".", atSymbolIndex) > atSymbolIndex + 1 &&
		email.indexOf(".") < email.length - 1
	);
}
