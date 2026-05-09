import { ADM_LAYER_LEVELS, getFillLayerIdForAdmLv, getOutlineLayerIdForAdmLv } from "./adm-map-layers";

export const USER_GEOMETRY_FILL_LAYER_ID = 'user-geometry-fill';
export const USER_GEOMETRY_OUTLINE_LAYER_ID = 'user-geometry-outline';

export const LAYERS_IDS_IN_ORDER = [
    'background',
    ...[...ADM_LAYER_LEVELS].reverse().map(getFillLayerIdForAdmLv),
    ...[...ADM_LAYER_LEVELS].reverse().map(getOutlineLayerIdForAdmLv),
]