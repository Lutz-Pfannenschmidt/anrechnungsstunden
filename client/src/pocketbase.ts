import PocketBase, { type OTPResponse } from "pocketbase";
import { Collections, type PartPointsResponse, type ResultsRecord, type TypedPocketBase, type UsersResponse } from "./pocketbase-types";

export const dev = process.env.NODE_ENV === "development";

// The URL of the PocketBase server MUST end with a slash
export const pbUrl = dev ? "http://localhost:8090/" : window.location.origin;

export const pb = new PocketBase(pbUrl) as TypedPocketBase;
pb.autoCancellation(false);

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
    if (path.startsWith("/admin/login")) {
        return
    }
    window.location.href = "/admin/login";
}

export function isLoggedIn(): boolean {
    return pb.authStore.isValid;
}

export async function refreshAuth() {
    if (pb.authStore.isValid && pb.authStore.record) {
        try {
            const auth = await pb
                .collection(pb.authStore.record.collectionName)
                .authRefresh();
            if (!auth || !auth.record) {
                logout();
                return false;
            }
        } catch (error) {
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

export async function putPartPoints(subject: string, grade: string, points: string, id?: string | null): Promise<PartPointsResponse | null> {
    let existsID = id || "";
    if (!existsID) {
        try {
            const record = await pb
                .collection(Collections.PartPoints)
                .getFirstListItem(`class="${subject}" && grade="${grade}"`);
            existsID = record.id || "";
        } catch (error) {
            existsID = "";
        }
    }
    if (existsID) {
        try {
            const record = await pb
                .collection(Collections.PartPoints)
                .update(existsID, { class: subject, grade: grade, points: points });
            return record;
        } catch (error) {
            return null;
        }
    } else {
        try {
            const record = await pb
                .collection(Collections.PartPoints)
                .create({ class: subject, grade: grade, points: points });
            return record;
        } catch (error) {
            return null;
        }
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