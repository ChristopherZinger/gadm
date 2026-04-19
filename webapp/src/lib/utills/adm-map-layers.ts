export const ADM_LAYER_LEVELS = [0, 1, 2, 3, 4, 5];

export function getOutlineLayerIdForAdmLv(lv: number): string {
	return `adm-${lv}-line`;
}

export function getFillLayerIdForAdmLv(lv: number): string {
	return `adm-${lv}-fill`;
}
