import { isLoggedIn, isSuperuser } from "@src/pocketbase";

const publicPaths = ["/login", "/admin/login"];

function checkRoute(path: string) {
    if (publicPaths.includes(path)) {
        if (isLoggedIn()) {
            if (isSuperuser()) {
                window.location.href = "/admin";
            } else {
                window.location.href = "/";
            }
            return;
        }
        return;
    }

    if (path === "/") {
        if (isLoggedIn()) {
            if (isSuperuser()) {
                window.location.href = "/admin";
            }
            return;
        }
        window.location.href = "/login";
        return;
    }

    if (path.startsWith("/admin")) {
        if (isLoggedIn()) {
            if (isSuperuser()) {
                return;
            }
            window.location.href = "/";
            return;
        }
        window.location.href = "/admin/login";
        return;
    }

    window.location.href = "/404";
    return;
}

let path = window.location.pathname;
path = path.endsWith('/') ? path.slice(0, -1) : path
path = path.length === 0 ? "/" : path
console.log("Checking route: ", path);
checkRoute(path);
