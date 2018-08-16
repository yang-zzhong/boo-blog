package controller

import (
	"boo-blog/model"
	"database/sql"
	helpers "github.com/yang-zzhong/go-helpers"
	httprouter "github.com/yang-zzhong/go-httprouter"
	m "github.com/yang-zzhong/go-model"
	. "github.com/yang-zzhong/go-querybuilder"
)

type User struct{ *Controller }

func (this *User) Find(req *httprouter.Request, p *helpers.P) {
	user := model.NewUser()
	if keyword := req.FormValue("keyword"); keyword != "" {
		user.Repo().Quote(func(repo *Builder) {
			repo.Where("name", LIKE, "%"+keyword+"%").
				Or().Where("phone_number", LIKE, "%"+keyword+"%").
				Or().Where("email_addr", LIKE, "%"+keyword+"%")
		})
	}
	page := req.FormInt("page")
	if page == 0 {
		page = 1
	}
	user.Repo().Page(int(page), 10)
	user.Repo().OrderBy("created_at", ASC)
	this.profiles(user, req, p, func(_ map[string]interface{}, _ *model.User) {})
}

func (this *User) AboutMe(req *httprouter.Request, p *helpers.P) {
	user := model.NewUser()
	userFollow := model.NewUserFollow()
	if p.Get("type") == "ifollowed" {
		userFollow.Repo().
			Select("followed").
			Where("user_id", p.Get("visitor_id"))
	} else if p.Get("type").(string) == "followedme" {
		userFollow.Repo().
			Select("user_id").
			Where("followed", p.Get("visitor_id"))
	} else {
		this.String("错误的类型", 500)
		return
	}
	user.Repo().WhereInQuery("id", userFollow.Repo().Builder)
	page := req.FormInt("page")
	if page == 0 {
		page = 1
	}
	user.Repo().Page(int(page), 10)
	this.profiles(user, req, p, func(_ map[string]interface{}, _ *model.User) {})
}

func (this *User) Profile(p *helpers.P) {
	user := model.NewUser()
	if i, exists, err := user.Repo().Find(p.Get("user_id")); err != nil {
		this.InternalError(err)
		return
	} else if !exists {
		this.String("用户未找到", 404)
		return
	} else {
		this.Json(i.(*model.User).Profile(), 200)
	}
}

func (this *User) One(req *httprouter.Request, p *helpers.P) {
	user := model.NewUser()
	user.Repo().With("current_theme").Where("name", p.Get("name"))
	if i, exist, err := user.Repo().One(); err != nil {
		this.InternalError(err)
		return
	} else if !exist {
		this.String("没有找到用户", 404)
		return
	} else {
		user = i.(*model.User)
		result := user.Map()
		if p.Get("visitor_id") != nil {
			userFollow := model.NewUserFollow()
			userFollow.Repo().
				Where("user_id", p.Get("visitor_id")).
				Where("followed", user.Id)
			if count, err := userFollow.Repo().Count(); err != nil {
				this.InternalError(err)
				return
			} else {
				result["i_followed"] = count > 0
			}
		} else {
			result["i_followed"] = false
		}
		if theme, err := i.(*model.User).One("current_theme"); err != nil {
			this.InternalError(err)
			return
		} else if theme != nil {
			result["theme"] = theme.(*model.Theme).Content
		}
		this.Json(result, 200)
	}
}

func (this *User) SaveUserInfo(req *httprouter.Request, p *helpers.P) {
	user := model.NewUser()
	if i, exist, err := user.Repo().Find(p.Get("visitor_id")); err != nil {
		this.InternalError(err)
	} else if !exist {
		this.String("用户不存在", 500)
	} else {
		user = i.(*model.User)
	}
	if req.FormValue("bio") != "" {
		user.Bio = req.FormValue("bio")
	}
	if req.FormValue("portrait_image_id") != "" {
		user.PortraitImageId = req.FormValue("portrait_image_id")
	}
	if req.FormValue("blog_name") != "" {
		user.BlogName = req.FormValue("blog_name")
	}
	if req.FormValue("name") != "" {
		user.Repo().Where("name", user.Name).Where("id", NEQ, p.Get("visitor_id"))
		if exists, err := user.Repo().Count(); err != nil {
			this.InternalError(err)
			return
		} else if exists > 0 {
			this.String("名字已存在", 500)
			return
		} else {
			user.Name = req.FormValue("name")
		}
	}
	if err := user.Save(); err != nil {
		this.InternalError(err)
	}
}

func (this *User) profiles(user *model.User, req *httprouter.Request, p *helpers.P, h handleItem) {
	user.Repo().With("current_theme")
	user.Repo().WithCustom("blogs", func(b interface{}) (data m.NexusValues, err error) {
		repo := b.(m.Model).Repo()
		repo.OrderBy("thumb_up", DESC).OrderBy("thumb_down", ASC)
		repo.Limit(3)
		if d, e := repo.Fetch(); e != nil {
			err = e
		} else {
			da := make(map[interface{}]interface{})
			for i, val := range d {
				da[i] = val
			}
			data = m.NewNexusValues(da)
		}
		return
	})
	if p.Get("visitor_id") != nil {
		user.Repo().WithCustom("followed", func(uf interface{}) (nv m.NexusValues, err error) {
			repo := uf.(m.Model).Repo()
			userIds := []uint32{}
			repo.Where("user_id", p.Get("visitor_id"))
			repo.Select("followed")
			err = repo.Query(func(rows *sql.Rows, _ []string) error {
				var userId uint32
				if err := rows.Scan(&userId); err != nil {
					return err
				}
				userIds = append(userIds, userId)

				return nil
			})
			if err == nil {
				nv = &userIdsIn{userIds}
			}
			return
		})
	}
	if ms, err := user.Repo().Fetch(); err != nil {
		this.InternalError(err)
	} else {
		result := []map[string]interface{}{}
		for _, i := range ms {
			profile := i.(*model.User).Profile()
			if p.Get("visitor_id") != nil {
				if followed, err := i.(*model.User).Many("followed"); err != nil {
					this.InternalError(err)
					return
				} else {
					profile["i_followed"] = followed
				}
			} else {
				profile["i_followed"] = false
			}
			if blogs, err := i.(*model.User).Many("blogs"); err != nil {
				this.InternalError(err)
				return
			} else {
				bs := []map[string]interface{}{}
				for _, b := range blogs.(map[interface{}]interface{}) {
					blog := b.(*model.Blog)
					bs = append(bs, map[string]interface{}{
						"url_id":     blog.UrlId,
						"title":      blog.Title,
						"overview":   blog.Overview,
						"created_at": blog.CreatedAt,
					})
				}
				profile["recommands"] = bs
			}
			h(profile, i.(*model.User))
			result = append(result, profile)
		}
		this.Json(result, 200)
	}
}

type handleItem func(profile map[string]interface{}, user *model.User)

type userIdsIn struct {
	userIds []uint32
}

func (uii *userIdsIn) DataOf(user interface{}, _ m.Nexus) interface{} {
	for _, userId := range uii.userIds {
		if user.(*model.User).Id == userId {
			return true
		}
	}
	return false
}
