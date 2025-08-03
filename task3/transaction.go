package task3

import "gorm.io/gorm"

type Account struct {
	ID        uint64 `gorm:"primary_key;auto_increment"`
	AccountNo string
	Balance   float64
}
type Transaction struct {
	ID            uint64 `gorm:"primary_key;auto_increment"`
	FromAccountID uint64
	FromAccount   Account `gorm:"foreignKey:FromAccountID"`
	ToAccountID   uint64
	ToAccount     Account `gorm:"foreignKey:ToAccountID"`
	Amount        float64
}

func Init() {
	db := ConnectDb()
	db.AutoMigrate(&Account{})
	db.AutoMigrate(&Transaction{})
}

func CreateData() {
	db := ConnectDb()
	var accounts = []Account{{AccountNo: "123", Balance: 88888888.88}, {AccountNo: "567", Balance: 88888888.88}}

	db.Debug().Create(&accounts)
}

func TransferMoney(fromAccountID, toAccountID uint64, amount float64) error {
	db := ConnectDb()
	tx := db.Begin()

	var fromAccount Account
	//检查账户 A 的余额是否足够
	if err := tx.First(&fromAccount, fromAccountID).Error; err != nil {
		tx.Rollback()
		return err
	}

	if fromAccount.Balance < amount {
		tx.Rollback()
		return gorm.ErrInvalidValue
	}

	// 更新账户余额
	if err := tx.Model(&Account{}).Where("id = ?", fromAccountID).
		Update("balance", gorm.Expr("balance - ?", amount)).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(&Account{}).Where("id = ?", toAccountID).
		Update("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
		tx.Rollback()
		return err
	}
	// 记录交易
	transaction := Transaction{
		FromAccountID: fromAccountID,
		ToAccountID:   toAccountID,
		Amount:        amount,
	}
	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
