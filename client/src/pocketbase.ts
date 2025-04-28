import PocketBase, { type OTPResponse } from "pocketbase";
import { Collections, type AcronymsResponse, type ResultsRecord, type TeacherDataResponse, type TypedPocketBase, type UsersResponse } from "./pocketbase-types";

export const dev = process.env.NODE_ENV === "development";

// The URL of the PocketBase server MUST end with a slash
export const pbUrl = dev ? "http://localhost:8090/" : window.location.origin;

export const pb = new PocketBase(pbUrl) as TypedPocketBase;
pb.autoCancellation(false);

// Sets the pocketbase instance to use the impersonated user for all requests and returns whether the impersonation was successful
export async function impersonate(userId: string, duration = 0): Promise<boolean> {
    try {
        const client = await pb.collection(Collections.Users).impersonate(userId, duration);

        const token = pb.authStore.token;
        const record = pb.authStore.record;

        localStorage.setItem("_superuser", JSON.stringify({ token, record }));
        localStorage.setItem("_superuser_route", window.location.pathname);

        pb.authStore.clear();
        pb.authStore.save(client.authStore.token, client.authStore.record);
        return true;
    } catch (error) {
        return false;
    }
}

export async function requestOTP(username: string, collection = "users"): Promise<OTPResponse | null> {
    try {
        const o = await pb.collection(collection).requestOTP(username);
        return o;
    } catch (error) {
        return null;
    }
}

export async function loginOTP(requestId: string, otp: string, collection = "users") {
    try {
        await pb.collection(collection).authWithOTP(requestId, otp);
    } catch (error) {
        return false;
    }
    return isLoggedIn();
}

export async function login(username: string, password: string, collection = "users") {
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

export function isLoggedIn(): boolean {
    return pb.authStore.isValid;
}

export async function refreshAuth() {
    if (pb.authStore.isValid && pb.authStore.record) {
        const auth = await pb
            .collection(pb.authStore.record.collectionName)
            .authRefresh();
        if (!auth || !auth.record) {
            logout();
            return false;
        }
    } else {
        logout();
        return false;
    }

    return isLoggedIn();
}

export function isSuperuser() {
    return pb.authStore.record?.collectionName === "_superusers";
}

// Create a new user
export async function createUser(email: string, username: string, short: string) {
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

export async function getUserId(email: string) {
    try {
        const record = await pb
            .collection("users")
            .getFirstListItem(`email="${email}"`);
        return record.id;
    } catch (error) {
        return "";
    }
}

export async function putData(uID: string, yearID: string, avgTime: number) {
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

export async function putExamPoints(subject: string, grade: string, points: number | string, id: string | null): Promise<string | null> {
    if (id) {
        try {
            const record = await pb
                .collection("exam_points")
                .update(id, { subject: subject, grade: grade, points: points });
            return record.id;
        } catch (error) {
            return null;
        }
    } else {
        try {
            const record = await pb
                .collection("exam_points")
                .create({ subject: subject, grade: grade, points: points });
            return record.id;
        } catch (error) {
            return null;
        }
    }
}

export async function putClassLead(teacherID: string, year: string | number, semester: string | number, points: number): Promise<string | null> {
    let existsID = "";
    try {
        const record = await pb
            .collection(Collections.ClassLead)
            .getFirstListItem(
                pb.filter("teacher={:teacherID} && year={:year} && semester={:semester}", {
                    teacherID,
                    year, semester
                })
            );
        existsID = record.id || "";
    } catch (error) {
        existsID = "";
    }

    if (existsID) {
        try {
            const record = await pb
                .collection(Collections.ClassLead)
                .update(existsID, { percentage: points });
            return record.id;
        } catch (error) {
            return null;
        }
    } else {
        try {
            const record = await pb
                .collection(Collections.ClassLead)
                .create({ teacher: teacherID, year: year, semester: semester, percentage: points });
            return record.id;
        } catch (error) {
            return null;
        }
    }
}

export async function getClassLead(teacherID: string, year: string | number, semester: string | number) {
    try {
        const record = await pb
            .collection(Collections.ClassLead)
            .getFirstListItem(
                pb.filter("teacher={:teacherID} && year={:year} && semester={:semester}", {
                    teacherID,
                    year, semester
                })
            );
        return record;
    } catch (error) {
        return null;
    }
}

export async function getTeacherData(teacherID: string, year: string | number, semester: string | number): Promise<TeacherDataResponse<unknown>[]> {
    try {
        const records = await pb.collection("teacher_data").getFullList({
            filter: pb.filter("teacher={:teacherID} && year={:year} && semester={:semester}", { teacherID, year, semester })
        });
        return records;
    } catch (error) {
        return [];
    }
}

export async function deleteYears(yearID: string) {
    try {
        await pb.collection("years").delete(yearID);

        return true;
    } catch (error) {
        return false;
    }
}

function genText(data: { [key: string]: number }): string {
    let text = "";
    for (const [i, [teacher, hours]] of Object.entries(data).entries() as Iterable<[number, [string, number]]>) {
        text += `${i + 10000};;;"${teacher.toUpperCase()}";"500","${hours.toFixed(3)};;;"LK";0.000;0.000;;${hours.toFixed(3)}\n`;
    }
    return text;
}

export async function putResults(semester: string, data: { [key: string]: number }, lead_points: number): Promise<ResultsRecord | null> {
    let existsID = "";
    try {
        const record = await pb
            .collection(Collections.Results)
            .getFirstListItem(`semester="${semester}"`);
        existsID = record.id || "";
    } catch (error) {
        existsID = "";
    }

    if (existsID) {
        try {
            const record = await pb
                .collection(Collections.Results)
                .update(existsID, { data: data, untis: genText(data), lead_points: lead_points });
            return record;
        } catch (error) {
            return null;
        }
    } else {
        try {
            const record = await pb
                .collection(Collections.Results)
                .create({ semester: semester, data: data, untis: genText(data), lead_points: lead_points });
            return record;
        } catch (error) {
            return null;
        }
    }
}

export function formatDateToDDMMYYYY(date: Date): string {
    const day = String(date.getDate()).padStart(2, "0");
    const month = String(date.getMonth() + 1).padStart(2, "0"); // Months are zero-based
    const year = date.getFullYear();

    return `${day}.${month}.${year}`;
}