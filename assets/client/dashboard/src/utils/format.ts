export function FormatToRaw(content: Record<string, string>): string {
  let result: string = "";
  Object.keys(content).forEach(function (key) {
    result += `${key}: ${content[key]}\n`;
  });

  return result;
}
