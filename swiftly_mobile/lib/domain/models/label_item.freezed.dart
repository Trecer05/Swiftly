// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'label_item.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

T _$identity<T>(T value) => value;

final _privateConstructorUsedError = UnsupportedError(
  'It seems like you constructed your class using `MyClass._()`. This constructor is only meant to be used by freezed and you are not supposed to need it nor use it.\nPlease check the documentation here for more information: https://github.com/rrousselGit/freezed#adding-getters-and-methods-to-our-models',
);

/// @nodoc
mixin _$LabelItem {
  String get id => throw _privateConstructorUsedError;
  String? get cardId => throw _privateConstructorUsedError;
  String? get userId => throw _privateConstructorUsedError;
  String get title => throw _privateConstructorUsedError;
  Color get color => throw _privateConstructorUsedError;

  /// Create a copy of LabelItem
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  $LabelItemCopyWith<LabelItem> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $LabelItemCopyWith<$Res> {
  factory $LabelItemCopyWith(LabelItem value, $Res Function(LabelItem) then) =
      _$LabelItemCopyWithImpl<$Res, LabelItem>;
  @useResult
  $Res call({
    String id,
    String? cardId,
    String? userId,
    String title,
    Color color,
  });
}

/// @nodoc
class _$LabelItemCopyWithImpl<$Res, $Val extends LabelItem>
    implements $LabelItemCopyWith<$Res> {
  _$LabelItemCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  /// Create a copy of LabelItem
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? id = null,
    Object? cardId = freezed,
    Object? userId = freezed,
    Object? title = null,
    Object? color = null,
  }) {
    return _then(
      _value.copyWith(
            id:
                null == id
                    ? _value.id
                    : id // ignore: cast_nullable_to_non_nullable
                        as String,
            cardId:
                freezed == cardId
                    ? _value.cardId
                    : cardId // ignore: cast_nullable_to_non_nullable
                        as String?,
            userId:
                freezed == userId
                    ? _value.userId
                    : userId // ignore: cast_nullable_to_non_nullable
                        as String?,
            title:
                null == title
                    ? _value.title
                    : title // ignore: cast_nullable_to_non_nullable
                        as String,
            color:
                null == color
                    ? _value.color
                    : color // ignore: cast_nullable_to_non_nullable
                        as Color,
          )
          as $Val,
    );
  }
}

/// @nodoc
abstract class _$$LabelItemImplCopyWith<$Res>
    implements $LabelItemCopyWith<$Res> {
  factory _$$LabelItemImplCopyWith(
    _$LabelItemImpl value,
    $Res Function(_$LabelItemImpl) then,
  ) = __$$LabelItemImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call({
    String id,
    String? cardId,
    String? userId,
    String title,
    Color color,
  });
}

/// @nodoc
class __$$LabelItemImplCopyWithImpl<$Res>
    extends _$LabelItemCopyWithImpl<$Res, _$LabelItemImpl>
    implements _$$LabelItemImplCopyWith<$Res> {
  __$$LabelItemImplCopyWithImpl(
    _$LabelItemImpl _value,
    $Res Function(_$LabelItemImpl) _then,
  ) : super(_value, _then);

  /// Create a copy of LabelItem
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? id = null,
    Object? cardId = freezed,
    Object? userId = freezed,
    Object? title = null,
    Object? color = null,
  }) {
    return _then(
      _$LabelItemImpl(
        id:
            null == id
                ? _value.id
                : id // ignore: cast_nullable_to_non_nullable
                    as String,
        cardId:
            freezed == cardId
                ? _value.cardId
                : cardId // ignore: cast_nullable_to_non_nullable
                    as String?,
        userId:
            freezed == userId
                ? _value.userId
                : userId // ignore: cast_nullable_to_non_nullable
                    as String?,
        title:
            null == title
                ? _value.title
                : title // ignore: cast_nullable_to_non_nullable
                    as String,
        color:
            null == color
                ? _value.color
                : color // ignore: cast_nullable_to_non_nullable
                    as Color,
      ),
    );
  }
}

/// @nodoc

class _$LabelItemImpl implements _LabelItem {
  const _$LabelItemImpl({
    required this.id,
    required this.cardId,
    required this.userId,
    required this.title,
    required this.color,
  });

  @override
  final String id;
  @override
  final String? cardId;
  @override
  final String? userId;
  @override
  final String title;
  @override
  final Color color;

  @override
  String toString() {
    return 'LabelItem(id: $id, cardId: $cardId, userId: $userId, title: $title, color: $color)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$LabelItemImpl &&
            (identical(other.id, id) || other.id == id) &&
            (identical(other.cardId, cardId) || other.cardId == cardId) &&
            (identical(other.userId, userId) || other.userId == userId) &&
            (identical(other.title, title) || other.title == title) &&
            (identical(other.color, color) || other.color == color));
  }

  @override
  int get hashCode =>
      Object.hash(runtimeType, id, cardId, userId, title, color);

  /// Create a copy of LabelItem
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  @pragma('vm:prefer-inline')
  _$$LabelItemImplCopyWith<_$LabelItemImpl> get copyWith =>
      __$$LabelItemImplCopyWithImpl<_$LabelItemImpl>(this, _$identity);
}

abstract class _LabelItem implements LabelItem {
  const factory _LabelItem({
    required final String id,
    required final String? cardId,
    required final String? userId,
    required final String title,
    required final Color color,
  }) = _$LabelItemImpl;

  @override
  String get id;
  @override
  String? get cardId;
  @override
  String? get userId;
  @override
  String get title;
  @override
  Color get color;

  /// Create a copy of LabelItem
  /// with the given fields replaced by the non-null parameter values.
  @override
  @JsonKey(includeFromJson: false, includeToJson: false)
  _$$LabelItemImplCopyWith<_$LabelItemImpl> get copyWith =>
      throw _privateConstructorUsedError;
}
