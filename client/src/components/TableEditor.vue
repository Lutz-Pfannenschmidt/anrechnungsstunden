<script setup lang="ts">
import { pb } from '@src/pocketbase';
import { Collections } from '@src/pocketbase-types';
import { customPrompt, toast } from '@src/toast';
import { ref } from 'vue'

const points = await pb.collection(Collections.PartPoints).getFullList();

const mapped = ref(new Map<string, Map<string, number>>());
const grades = ref<string[]>([]);

for (const point of points) {
    if (!mapped.value.has(point.class)) {
        mapped.value.set(point.class, new Map());
    }
    mapped.value.get(point.class)?.set(point.grade, point.points);
    if (!grades.value.includes(point.grade)) {
        grades.value.push(point.grade);
    }
    grades.value.sort()
}


pb.collection(Collections.PartPoints).subscribe('*', (e) => {
    const item = e.record;
    if (e.action === 'create' || e.action === 'update') {
        if (!mapped.value.has(item.class)) {
            mapped.value.set(item.class, new Map());
        }
        mapped.value.get(item.class)?.set(item.grade, item.points);
    } else if (e.action === 'delete') {
        mapped.value.get(item.class)?.delete(item.grade);
    }

    if (!grades.value.includes(item.grade)) {
        grades.value.push(item.grade);
    } else if (mapped.value.size === 0) {
        grades.value = [];
    }
    grades.value.sort()
});

async function addClass() {
    const className = await customPrompt("Fachabkürzung wie in Leistungsdaten (z.B. 'D', 'M', ...)", "text")
    if (!className) return;

    if (mapped.value.has(className)) {
        toast("warning", "Diese Fachabkürzung existiert bereits!");
        return;
    }

    const newClassPoints = new Map<string, number>();
    for (const grade of grades.value) {
        newClassPoints.set(grade, 0);
    }
    mapped.value.set(className, newClassPoints);
}

async function addGrade() {
    const grade = await customPrompt("Neue Stufe hinzufügen (z.B. '10', 'EF', ...)", "text")
    if (!grade) return;

    if (grades.value.includes(grade)) {
        toast("warning", "Diese Stufe existiert bereits!");
        return;
    }

    grades.value.push(grade);
    grades.value.sort();

}

function changeName(className: string, newName: string) {
    if (newName === className || !newName) return;

    const classPoints = mapped.value.get(className);
    if (!classPoints) return;

    mapped.value.delete(className);
    mapped.value.set(newName, classPoints);

    console.log(`Renamed class ${className} to ${newName}`);
}

function delName(className: string) {
    if (!mapped.value.has(className)) return;

    const classPoints = mapped.value.get(className);
    if (!classPoints) return;

    mapped.value.delete(className);
    for (const grade of classPoints.keys()) {
        updateOrDeletePoint(className, grade, null);

    }

    console.log(`Deleted class ${className}`);

}

async function updateOrDeletePoint(className: string, grade: string, points: number | null) {
    if (points === null) {
        await pb.collection(Collections.PartPoints).delete(`${className}-${grade}`);
    } else {
        await pb.collection(Collections.PartPoints).update(`${className}-${grade}`, {
            class: className,
            grade: grade,
            points: points
        });
    }
}

async function createOrUpdatePoint(className: string, grade: string, points: number) {
    if (mapped.value.has(className) && mapped.value.get(className)!.has(grade)) {
        await updateOrDeletePoint(className, grade, points);
    } else {
        await pb.collection(Collections.PartPoints).create({
            class: className,
            grade: grade,
            points: points,
            id: `${className}-${grade}`
        });
    }
}

function changePoints(className: string, grade: string, points: string) {
    const p = parseFloat(points.replaceAll('.', ' ').replaceAll(',', '.'));
    if (isNaN(p)) {
        toast("error", "Bitte eine gültige Zahl eingeben!");
        return;
    }

    createOrUpdatePoint(className, grade, p);
}

</script>

<template>
    <div class="overflow-x-auto table-fixed rounded-box border border-base-content/30">
        <table class="table" :class="grades.length > 2 && 'table-zebra'">
            <thead>
                <tr class=" border-b border-base-content/30">
                    <th class="border-r border-base-content/30">/</th>
                    <th class="border-r border-base-content/30" v-for="className in mapped.keys()" :key="className">
                        <div class="flex items-center gap-4">
                            <input type="text" class="input grow"
                                @change="changeName(className, ($event.target as HTMLInputElement).value)"
                                :value="className"></input>
                            <button @click="delName(className)" class="btn btn-outline btn-error">Löschen</button>
                        </div>
                    </th>
                    <th class="w-0 border-b border-transparent"><button @click="addClass"
                            class="btn btn-outline btn-success">+</button>
                    </th>
                </tr>
            </thead>
            <tbody class="border border-base-content/5">
                <tr v-for="grade in grades" :key="grade" class="border-t border-base-content/30">
                    <th class="border-r border-base-content/30 w-1">{{ grade }}</th>
                    <td class="border-r border-base-content/30" v-for="className in mapped.keys()" :key="className"
                        :class="mapped.has(className) && mapped.get(className)?.has(grade) ? 'text-base-content' : 'bg-gray-300'">
                        <input type="text"
                            @change="changePoints(className, grade, ($event.target as HTMLInputElement).value)">{{
                                mapped.get(className)?.get(grade) ?? '-' }}</input>
                    </td>
                    <td class="border-t border-transparent"></td>
                </tr>

                <tr class="border-t border-base-content/30">
                    <th class="w-0"><button @click="addGrade" class="btn btn-outline btn-success">+</button></th>
                    <td v-for="className in mapped.keys()" :key="className"></td>
                    <td class="border-t border-transparent"></td>
                </tr>


                <tr v-if="mapped.size === 0">
                    <td colspan="100%" class="text-center">
                        <span class="text-base-content font-bold text-lg">Keine Daten vorhanden!</span>
                    </td>
                </tr>
            </tbody>

        </table>

    </div>
</template>