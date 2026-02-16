export const NOM35_RISK_LEVELS = ["nulo", "bajo", "medio", "alto", "muy_alto"] as const;

export type Nom35RiskLevel = (typeof NOM35_RISK_LEVELS)[number];

const RISK_LEVEL_SET = new Set<string>(NOM35_RISK_LEVELS);

export const NOM35_RISK_LABELS: Record<Nom35RiskLevel, string> = {
  nulo: "Nulo",
  bajo: "Bajo",
  medio: "Medio",
  alto: "Alto",
  muy_alto: "Muy Alto",
};

export const NOM35_RISK_HEX_COLORS: Record<Nom35RiskLevel, string> = {
  nulo: "#16a34a",
  bajo: "#0284c7",
  medio: "#d97706",
  alto: "#ea580c",
  muy_alto: "#dc2626",
};

export const NOM35_RISK_BADGE_CLASSES: Record<Nom35RiskLevel, string> = {
  nulo: "risk-level-nulo border",
  bajo: "risk-level-bajo border",
  medio: "risk-level-medio border",
  alto: "risk-level-alto border",
  muy_alto: "risk-level-muy-alto border",
};

/**
 * Guide III risk thresholds:
 * - nulo: 0-49
 * - bajo: 50-74
 * - medio: 75-98
 * - alto: 99-139
 * - muy_alto: 140+
 */
export function getGuideIIIRiskLevelByScore(score: number | null): Nom35RiskLevel {
  if (score === null || score < 50) return "nulo";
  if (score < 75) return "bajo";
  if (score < 99) return "medio";
  if (score < 140) return "alto";
  return "muy_alto";
}

export function isNom35RiskLevel(value: string): value is Nom35RiskLevel {
  return RISK_LEVEL_SET.has(value);
}

export function getNom35RiskHexColor(value: string, fallback = "#6b7280"): string {
  return isNom35RiskLevel(value) ? NOM35_RISK_HEX_COLORS[value] : fallback;
}

export function getNom35RiskLabel(value: string): string {
  return isNom35RiskLevel(value) ? NOM35_RISK_LABELS[value] : value;
}

export function getNom35RiskBadgeClass(value: string): string {
  return isNom35RiskLevel(value) ? NOM35_RISK_BADGE_CLASSES[value] : "bg-muted text-muted-foreground border border-border";
}
