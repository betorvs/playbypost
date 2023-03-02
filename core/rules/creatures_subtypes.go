package rules

import "github.com/betorvs/playbypost/core/types"

type CreatureSubtype int

const (
	Air CreatureSubtype = iota
	Angel
	AquaticCreatureSubtype
	Archon
	Augmented
	Chaotic
	Cold
	Demon
	Devil
	Earth
	Evil
	Extraplanar
	Fire
	GoodCreatureSubtype
	IncorporealSubtype
	Lawful
	Native
	Psionic
	Shapechanger
	Swarm
	Water
)

func (c CreatureSubtype) String() string {
	switch c {
	case Air:
		return "air"
	case Angel:
		return "angel"
	case AquaticCreatureSubtype:
		return "aquatic"
	case Archon:
		return "archon"
	case Augmented:
		return "augmented"
	case Chaotic:
		return "chaotic"
	case Cold:
		return "cold"
	case Demon:
		return "demon"
	case Devil:
		return "devil"
	case Earth:
		return "earth"
	case Evil:
		return "evil"
	case Extraplanar:
		return "extraplanar"
	case Fire:
		return "fire"
	case GoodCreatureSubtype:
		return "good"
	case IncorporealSubtype:
		return "incorporeal"
	case Lawful:
		return "lawful"
	case Native:
		return "native"
	case Psionic:
		return "psionic"
	case Shapechanger:
		return "shapechanger"
	case Swarm:
		return "swarm"
	case Water:
		return "water"
	}
	return types.Unknown
}

type HumanoidSubtype int

const (
	AquaticHumanoidSubType HumanoidSubtype = iota
	Dwarf
	Elf
	Gnoll
	Gnome
	Goblinoid
	Halfling
	Human
	Orc
	Reptilian
)

func (h HumanoidSubtype) String() string {
	switch h {
	case AquaticHumanoidSubType:
		return "aquatic"
	case Dwarf:
		return "dwarf"
	case Elf:
		return "elf"
	case Gnoll:
		return "gnoll"
	case Gnome:
		return "gnome"
	case Goblinoid:
		return "goblinoid"
	case Halfling:
		return "halfling"
	case Human:
		return "human"
	case Orc:
		return "orc"
	case Reptilian:
		return "reptilian"
	}
	return types.Unknown
}

type EnemyGroups int

const (
	AberrationEnemyGroups EnemyGroups = iota
	AnimalEnemyGroups
	ConstructEnemyGroups
	DragonEnemyGroups
	ElementalEnemyGroups
	FeyEnemyGroups
	GiantEnemyGroups
	HumanoidAquaticEnemyGroups
	HumanoidDwarfEnemyGroups
	HumanoidElfEnemyGroups
	HumanoidGoblinoidEnemyGroups
	HumanoidGnollEnemyGroups
	HumanoidGnomeEnemyGroups
	HumanoidHalflingEnemyGroups
	HumanoidHumanEnemyGroups
	HumanoidOrcEnemyGroups
	HumanoidReptilianEnemyGroups
	MagicalBeastEnemyGroups
	MonstrousHumanoidEnemyGroups
	OozeEnemyGroups
	OutsiderAirEnemyGroups
	OutsiderChaoticEnemyGroups
	OutsiderEarthEnemyGroups
	OutsiderEvilEnemyGroups
	OutsiderFireEnemyGroups
	OutsiderGoodEnemyGroups
	OutsiderLawfulEnemyGroups
	OutsiderNativeEnemyGroups
	OutsiderWaterEnemyGroups
	PlantEnemyGroups
	UndeadEnemyGroups
	VerminEnemyGroups
)

func (e EnemyGroups) String() string {
	switch e {
	case AberrationEnemyGroups:
		return "aberration"
	case AnimalEnemyGroups:
		return "animal"
	case ConstructEnemyGroups:
		return "construct"
	case DragonEnemyGroups:
		return "dragon"
	case ElementalEnemyGroups:
		return "elemental"
	case FeyEnemyGroups:
		return "fey"
	case GiantEnemyGroups:
		return "giant"
	case HumanoidAquaticEnemyGroups:
		return "humanoid aquatic"
	case HumanoidDwarfEnemyGroups:
		return "humanoid dwarf"
	case HumanoidElfEnemyGroups:
		return "humanoid elf"
	case HumanoidGoblinoidEnemyGroups:
		return "humanoid goblinoid"
	case HumanoidGnollEnemyGroups:
		return "humanoid gnoll"
	case HumanoidGnomeEnemyGroups:
		return "humanoid gnome"
	case HumanoidHalflingEnemyGroups:
		return "humanoid halfling"
	case HumanoidHumanEnemyGroups:
		return "humanoid human"
	case HumanoidOrcEnemyGroups:
		return "humanoid orc"
	case HumanoidReptilianEnemyGroups:
		return "humanoid reptilian"
	case MagicalBeastEnemyGroups:
		return "magical beast"
	case MonstrousHumanoidEnemyGroups:
		return "monstrous humanoid"
	case OozeEnemyGroups:
		return "ooze"
	case OutsiderAirEnemyGroups:
		return "outsider air"
	case OutsiderChaoticEnemyGroups:
		return "outsider chaotic"
	case OutsiderEarthEnemyGroups:
		return "outsider earth"
	case OutsiderEvilEnemyGroups:
		return "outsider evil"
	case OutsiderFireEnemyGroups:
		return "outsider fire"
	case OutsiderGoodEnemyGroups:
		return "outsider good"
	case OutsiderLawfulEnemyGroups:
		return "outsider lawful"
	case OutsiderNativeEnemyGroups:
		return "outsider native"
	case OutsiderWaterEnemyGroups:
		return "outsider water"
	case PlantEnemyGroups:
		return "plant"
	case UndeadEnemyGroups:
		return "undead"
	case VerminEnemyGroups:
		return "vermin"
	}
	return types.Unknown
}
