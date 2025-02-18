/**
 * Capitalize the first letter of a string
 *
 * @export
 * @param {string} string The string to capitalize
 * @returns {string} The string with the first letter capitalized
 */
export function toUpperCaseFirst(string) {
	return string.charAt(0).toUpperCase() + string.slice(1);
}
