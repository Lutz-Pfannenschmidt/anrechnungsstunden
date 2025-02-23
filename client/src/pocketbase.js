import PocketBase from "pocketbase";

export const dev = process.env.NODE_ENV === "development";

// The URL of the PocketBase server MUST end with a slash
export const pbUrl = dev ? "http://localhost:8090/" : window.location.origin;

export const pb = new PocketBase(pbUrl);
pb.autoCancellation(false);

export async function requestOTP(username, collection = "users") {
	try {
		const o = await pb.collection(collection).requestOTP(username);
		return o;
	} catch (error) {
		return null;
	}
}

export async function loginOTP(requestId, otp, collection = "users") {
	try {
		await pb.collection(collection).authWithOTP(requestId, otp);
	} catch (error) {
		return false;
	}
	return isLoggedIn();
}

export async function login(username, password, collection = "users") {
	try {
		await pb.collection(collection).authWithPassword(username, password);
	} catch (error) {
		return false;
	}
	return isLoggedIn();
}

/**
 * Logout the current user
 */
export function logout() {
	pb.authStore.clear();
	const path = window.location.pathname;
	if (!path.startsWith("/login") && !path.startsWith("/admin/login")) {
		window.location.href = "/login";
	}
}

/**
 * Check if a user is logged in
 * @returns {boolean} - Whether a user is logged in or not
 */
export function isLoggedIn() {
	return pb.authStore.isValid;
}

/**
 * Redirect to the login page if the user is not logged in
 */
export function mustBeLoggedIn() {
	if (!isLoggedIn()) {
		window.location.href = "/login";
	}
}

export function isSuperuser() {
	return pb.authStore.record?.collectionName === "_superusers";
}

/**
 * Create a new user
 * @param {string} email - The email of the user
 * @param {string} username - The username of the user
 * @param {string} short - The short name of the user
 * @returns {Promise<string>} - The ID of the created user
 */
export async function createUser(email, username, short) {
	const pwd = `${Math.random().toString(36).slice(-8)}Aa1!`;
	const data = {
		password: pwd,
		passwordConfirm: pwd,
		email: email.toLowerCase(),
		verified: true,
		name: username,
		short: short,
	};

	try {
		const record = await pb.collection("users").create(data);
		return record.id;
	} catch (error) {
		return "";
	}
}

export async function getUserId(email) {
	try {
		const record = await pb
			.collection("users")
			.getFirstListItem(`email="${email}"`);
		return record.id;
	} catch (error) {
		return "";
	}
}

export async function putData(uID, yearID, avgTime) {
	let existsID = "";
	try {
		const record = await pb
			.collection("time_data")
			.getFirstListItem(`user="${uID}" && semester="${yearID}"`);
		existsID = record.id || "";
	} catch (error) {
		existsID = "";
	}
	if (existsID) {
		try {
			const record = await pb
				.collection("time_data")
				.update(existsID, { avg_time: avgTime });
			return record.id;
		} catch (error) {
			return "";
		}
	} else {
		try {
			const record = await pb
				.collection("time_data")
				.create({ user: uID, semester: yearID, avg_time: avgTime });
			return record.id;
		} catch (error) {
			return "";
		}
	}
}

export async function putExamPoints(subject, grade, points, id) {
	if (id) {
		try {
			const record = await pb
				.collection("exam_points")
				.update(id, { subject: subject, grade: grade, points: points });
			return record.id;
		} catch (error) {
			return "";
		}
	} else {
		try {
			const record = await pb
				.collection("exam_points")
				.create({ subject: subject, grade: grade, points: points });
			return record.id;
		} catch (error) {
			return "";
		}
	}
}

export async function putClassLead(teacherID, year, semester, percentage) {
	let existsID = "";
	try {
		const record = await pb
			.collection("class_lead")
			.getFirstListItem(
				`teacher="${teacherID}" && year="${year}" && semester="${semester}"`,
			);
		existsID = record.id || "";
	} catch (error) {
		existsID = "";
	}
	if (existsID) {
		try {
			const record = await pb
				.collection("class_lead")
				.update(existsID, { percentage: percentage });
			return record.id;
		} catch (error) {
			return "";
		}
	} else {
		try {
			const record = await pb.collection("class_lead").create({
				teacher: teacherID,
				year: year,
				semester: semester,
				percentage: percentage,
			});
			return record.id;
		} catch (error) {
			return "";
		}
	}
}

export async function putTeacherData(
	teacherID,
	year,
	semester,
	grade,
	subject,
	students,
) {
	let existsID = "";
	try {
		const record = await pb
			.collection("teacher_data")
			.getFirstListItem(
				`teacher="${teacherID}" && year="${year}" && semester="${semester}" && subject="${subject}" && grade="${grade}"`,
			);
		existsID = record.id || "";
	} catch (error) {
		existsID = "";
	}
	if (existsID) {
		try {
			const record = await pb
				.collection("teacher_data")
				.update(existsID, { students: students });
			return record.id;
		} catch (error) {
			return "";
		}
	} else {
		try {
			const record = await pb.collection("teacher_data").create({
				teacher: teacherID,
				year: year,
				semester: semester,
				students: students,
				grade: grade,
				subject: subject,
			});
			return record.id;
		} catch (error) {
			return "";
		}
	}
}

export async function getTeacherData(teacherID, year, semester) {
	try {
		const records = await pb.collection("teacher_data").getList(1, 999, {
			filter: `teacher="${teacherID}" && year="${year}" && semester="${semester}"`,
		});
		return records;
	} catch (error) {
		return "";
	}
}

export async function deleteYears(yearID) {
	try {
		await pb.collection("years").delete(yearID);

		return true;
	} catch (error) {
		return false;
	}
}

export async function putResults(semester, data, lead_points) {
	let existsID = "";
	try {
		const record = await pb
			.collection("results")
			.getFirstListItem(`semester="${semester}"`);
		existsID = record.id || "";
	} catch (error) {
		existsID = "";
	}

	if (existsID) {
		try {
			const record = await pb
				.collection("results")
				.update(existsID, { data: data, lead_points: lead_points });
			return record;
		} catch (error) {
			return "";
		}
	} else {
		try {
			const record = await pb
				.collection("results")
				.create({ semester: semester, data: data, lead_points: lead_points });
			return record;
		} catch (error) {
			return "";
		}
	}
}
