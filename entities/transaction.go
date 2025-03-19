package entities

type Transaction struct {
	ID          int     `json:"id"`
	Amount      float64 `json:"amount" db:"jumlah"`         // Cocokkan dengan `jumlah`
	Category    string  `json:"category" db:"kategori"`     // Cocokkan dengan `kategori`
	Date        string  `json:"date" db:"tanggal"`          // Cocokkan dengan `tanggal`
	Description string  `json:"description" db:"deskripsi"` // Cocokkan dengan `deskripsi`

}
