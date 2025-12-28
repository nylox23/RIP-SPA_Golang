package repository

import (
	"fmt"
	"strings"
)

type Repository struct {
}

func NewRepository() (*Repository, error) {
	return &Repository{}, nil
}

type Acid struct {
	ID      int
	Name    string
	NameExt string
	Desc    string
	Hplus   int
	ImgURL  string
	Amount  int
}

func (r *Repository) GetAcids() ([]Acid, error) {
	acids := []Acid{
		{
			ID:      1,
			Name:    "HCl",
			NameExt: "Соляная кислота",
			Desc:    "Cильная одноосновная кислота, широко используемая в химической промышленности, металлургии и лабораторной практике. Она хорошо растворяет металлы и карбонаты, выделяя водород или углекислый газ. В быту встречается как компонент средств для очистки от накипи и ржавчины.",
			Hplus:   1,
			ImgURL:  "http://localhost:9000/acids/HCl.jpeg",
			Amount:  0,
		},
		{
			ID:      2,
			Name:    "H₂SO₄",
			NameExt: "Серная кислота",
			Desc:    "Cильная двухосновная кислота, одна из важнейших в промышленности («кровь химии»). Используется при производстве удобрений, красителей, аккумуляторов и взрывчатых веществ. Обладает обезвоживающими свойствами и активно взаимодействует с органическими веществами.",
			Hplus:   2,
			ImgURL:  "http://localhost:9000/acids/H2SO4.jpeg",
			Amount:  0,
		},
		{
			ID:      3,
			Name:    "H₃PO₄",
			NameExt: "Фосфорная кислота",
			Desc:    "Трёхосновная кислота, применяемая при производстве минеральных удобрений, антикоррозийных покрытий и в пищевой промышленности (пищевая добавка E338). Она слабее сильных минеральных кислот, но хорошо взаимодействует с основаниями и карбонатами. В напитках типа колы придаёт кисловатый вкус.",
			Hplus:   3,
			ImgURL:  "http://localhost:9000/acids/H3PO4.jpg",
			Amount:  0,
		},
		{
			ID:      4,
			Name:    "CH₃COOH",
			NameExt: "Уксусная кислота",
			Desc:    "Слабая органическая кислота, хорошо известная как главный компонент столового уксуса. Она образуется при брожении сахаров и применяется в пищевой промышленности, медицине и в качестве растворителя. Несмотря на слабую диссоциацию, при избытке реагирует с карбонатами с выделением CO₂.",
			Hplus:   1,
			ImgURL:  "http://localhost:9000/acids/CH3COOH.jpeg",
			Amount:  0,
		},
		{
			ID:      5,
			Name:    "H₂C₂O₄",
			NameExt: "Щавелевая кислота",
			Desc:    "Органическая двухосновная кислота, встречающаяся в щавеле, шпинате и ревене. Используется как отбеливающее и очищающее средство, а также в аналитической химии для титрования. В высоких дозах токсична, но в малых количествах естественно содержится в пище.",
			Hplus:   2,
			ImgURL:  "http://localhost:9000/acids/H2C2O4.jpeg",
			Amount:  0,
		},
	}
	if len(acids) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	return acids, nil
}

func (r *Repository) GetSelected() (map[Acid]float32, error) {
	var acids map[Acid]float32
	acids = make(map[Acid]float32)

	acids[Acid{
		ID:      1,
		Name:    "HCl",
		NameExt: "Соляная кислота",
		Desc:    "Cильная одноосновная кислота, широко используемая в химической промышленности, металлургии и лабораторной практике. Она хорошо растворяет металлы и карбонаты, выделяя водород или углекислый газ. В быту встречается как компонент средств для очистки от накипи и ржавчины.",
		Hplus:   1,
		ImgURL:  "http://localhost:9000/acids/HCl.jpeg",
		Amount:  10,
	}] = 3.07

	acids[Acid{
		ID:      3,
		Name:    "H₃PO₄",
		NameExt: "Фосфорная кислота",
		Desc:    "Трёхосновная кислота, применяемая при производстве минеральных удобрений, антикоррозийных покрытий и в пищевой промышленности (пищевая добавка E338). Она слабее сильных минеральных кислот, но хорошо взаимодействует с основаниями и карбонатами. В напитках типа колы придаёт кисловатый вкус.",
		Hplus:   3,
		ImgURL:  "http://localhost:9000/acids/H3PO4.jpg",
		Amount:  30,
	}] = 10.28

	acids[Acid{
		ID:      5,
		Name:    "H₂C₂O₄",
		NameExt: "Щавелевая кислота",
		Desc:    "Органическая двухосновная кислота, встречающаяся в щавеле, шпинате и ревене. Используется как отбеливающее и очищающее средство, а также в аналитической химии для титрования. В высоких дозах токсична, но в малых количествах естественно содержится в пище.",
		Hplus:   2,
		ImgURL:  "http://localhost:9000/acids/H2C2O4.jpeg",
		Amount:  20,
	}] = 4.97

	if len(acids) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	return acids, nil
}

func (r *Repository) GetAcid(id int) (Acid, error) {
	acids, err := r.GetAcids()
	if err != nil {
		return Acid{}, err
	}

	for _, acid := range acids {
		if acid.ID == id {
			return acid, nil
		}
	}
	return Acid{}, fmt.Errorf("кислота не найден")
}

func (r *Repository) GetAcidsByName(name string) ([]Acid, error) {
	acids, err := r.GetAcids()
	if err != nil {
		return []Acid{}, err
	}

	var result []Acid
	for _, acid := range acids {
		if strings.Contains(strings.ToLower(acid.NameExt), strings.ToLower(name)) {
			result = append(result, acid)
		}
	}

	return result, nil
}
