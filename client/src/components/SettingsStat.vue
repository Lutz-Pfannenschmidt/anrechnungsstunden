<script setup lang="ts">
import { pb } from '@src/pocketbase';
import { ref } from 'vue';
import Alert from './Alert.vue';

const settings = ref(await pb.settings.getAll() as Settings);

function getSettings() {
    pb.settings.getAll().then((data) => {
        settings.value = data as Settings;
    }).catch((error) => {
        console.error("Error fetching settings:", error);
    });
}

getSettings();
setInterval(() => {
    getSettings();
}, 5000);

interface Settings {
    smtp: {
        enabled: boolean;
        port: number;
        host: string;
        username: string;
        authMethod: string;
        tls: boolean;
        localName: string;
    };
    backups: {
        cron: string;
        cronMaxKeep: number;
        s3: {
            enabled: boolean;
            bucket: string;
            region: string;
            endpoint: string;
            accessKey: string;
            forcePathStyle: boolean;
        };
    };
    s3: {
        enabled: boolean;
        bucket: string;
        region: string;
        endpoint: string;
        accessKey: string;
        forcePathStyle: boolean;
    };
    meta: {
        appName: string;
        appURL: string;
        senderName: string;
        senderAddress: string;
        hideControls: boolean;
    };
    rateLimits: {
        rules: {
            label: string;
            audience: string;
            duration: number;
            maxRequests: number;
        }[];
        enabled: boolean;
    };
    trustedProxy: {
        headers: string[];
        useLeftmostIP: boolean;
    };
    batch: {
        enabled: boolean;
        maxRequests: number;
        timeout: number;
        maxBodySize: number;
    };
    logs: {
        maxDays: number;
        minLevel: number;
        logIP: boolean;
        logAuthId: boolean;
    };
}

async function setReccomended() {
    settings.value.batch.enabled = true;
    settings.value.batch.maxRequests = 1000;
    settings.value.batch.maxBodySize = 512000000;
    settings.value.batch.timeout = 30;
    settings.value.smtp.enabled = true;

    try {
        await pb.settings.update(settings.value);
        console.log("Settings updated to recommended values.");
    } catch (error) {
        console.error("Error updating settings:", error);
    }

    getSettings();
}
</script>

<template>
    <Alert v-if="!settings.batch.enabled" type="warning" text="Batch API ist ausgechaltet!"></Alert>
    <Alert v-if="!settings.smtp.enabled" type="warning"
        text="SMTP ist nicht aktiviert und es können keine Emails versendet werden! (Bitte in den Email Einstellungen aktivieren.)">
    </Alert>
    <Alert v-else-if="settings.smtp.host == 'smtp.example.com'" type="error"
        text="SMTP ist nicht konfiguriert! (Bitte in den Email Einstellungen konfigurieren.)">
    </Alert>
    <Alert v-if="settings.batch.maxRequests < 1000" type="warning"
        text="Batch API ist auf weniger als 1000 Anfragen begrenzt! (Bitte in den Batch Einstellungen konfigurieren.)">
    </Alert>
    <Alert v-if="!settings.batch.enabled ||
        settings.batch.maxRequests < 1000 === true" type="info"
        text="Nicht alle Einstellungen sind nach den empfohlenen Werten konfiguriert. (Einstellungen können unter dem Menüpunkt 'Datenbank Einstellungen' geändert werden.)">
        <button @click="setReccomended()" class="btn btn-sm btn-primary">Alles auf Standartwerte setzen (Nur bei
            Erstanwendung empfohlen)</button>
    </Alert>
</template>