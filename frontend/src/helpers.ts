function randomInteger(max: number) {
  return Math.floor(Math.random() * (max + 1))
}

export function randomHexColor() {
  const hr = randomInteger(255).toString(16).padStart(2, "0")
  const hg = randomInteger(255).toString(16).padStart(2, "0")
  const hb = randomInteger(255).toString(16).padStart(2, "0")

  return "#" + hr + hg + hb
}
