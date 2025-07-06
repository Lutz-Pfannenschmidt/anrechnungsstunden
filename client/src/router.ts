import { isLoggedIn, isSuperuser, pb } from "@src/pocketbase";

function checkRoute(path: string) {
    if (path === "/login") {
        if (isLoggedIn() && isSuperuser()) {
            window.location.href = "/admin";
            return;
        }
        pb.authStore.clear();
        window.location.href = "/admin/login";
        return;
    }

    if (path === "/admin/login") {
        if (isLoggedIn()) {
            if (isSuperuser()) {
                window.location.href = "/admin";
            } else {
                pb.authStore.clear();
            }
        }
        return;
    }

    if (path === "/") {
        if (isLoggedIn() && isSuperuser()) {
            window.location.href = "/admin";
            return;
        }
        pb.authStore.clear();
        window.location.href = "/admin/login";
        return;
    }

    if (path.startsWith("/admin")) {
        if (isLoggedIn() && isSuperuser()) {
            return;
        }
        pb.authStore.clear();
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
