import { ADM_LAYER_LEVELS, getFillLayerIdForAdmLv, getOutlineLayerIdForAdmLv } from "./adm-map-layers";

export const LAYERS_IDS_IN_ORDER = [
    'background',
    ...[...ADM_LAYER_LEVELS].reverse().map(getFillLayerIdForAdmLv),
    ...[...ADM_LAYER_LEVELS].reverse().map(getOutlineLayerIdForAdmLv)
]