import PocketBase, { type OTPResponse } from "pocketbase";
import { Collections, type ResultsRecord, type TeacherDataResponse, type TypedPocketBase } from "./pocketbase-types";

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

/**
 * Check if a user is logged in
 * @returns {boolean} - Whether a user is logged in or not
 */
export function isLoggedIn() {
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

    return pb.authStore.isValid;
}

/**
 * Redirect to the login page if the user is not logged in
 */
export async function mustBeLoggedIn() {
    if (pb.authStore.isValid && pb.authStore.record) {
        const auth = await pb
            .collection(pb.authStore.record.collectionName)
            .authRefresh();
        if (!auth || !auth.record) {
            window.location.href = "/login";
        }
    } else {
        window.location.href = "/login";
    }
}

/**
 * Redirect to the login page if the user is not logged in
 */
export async function mustNotBeLoggedIn() {
    if (pb.authStore.isValid && pb.authStore.record) {
        try {
            const auth = await pb
                .collection(pb.authStore.record.collectionName)
                .authRefresh();
            if (auth?.record) {
                switch (pb.authStore.record.collectionName) {
                    case "users":
                        window.location.href = "/";
                        break;
                    case "_superusers":
                        window.location.href = "/admin";
                        break;
                    default:
                        window.location.href = "/login";
                }
            }
        } catch {
            switch (pb.authStore.record.collectionName) {
                case "users":
                    window.location.href = "/";
                    break;
                case "_superusers":
                    window.location.href = "/admin";
                    break;
                default:
                    window.location.href = "/login";
            }
        }
    }
}

export function redirectLoggedIn() {
    switch (pb.authStore.record?.collectionName) {
        case "users":
            window.location.href = "/";
            break;
        case "_superusers":
            window.location.href = "/admin/";
            break;
        default:
            window.location.href = "/login/";
    }
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
                .update(existsID, { points: points });
            return record.id;
        } catch (error) {
            return null;
        }
    } else {
        try {
            const record = await pb
                .collection(Collections.ClassLead)
                .create({ teacher: teacherID, year: year, semester: semester, points: points });
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
                .update(existsID, { data: data, lead_points: lead_points });
            return record;
        } catch (error) {
            return null;
        }
    } else {
        try {
            const record = await pb
                .collection(Collections.Results)
                .create({ semester: semester, data: data, lead_points: lead_points });
            return record;
        } catch (error) {
            return null;
        }
    }
}